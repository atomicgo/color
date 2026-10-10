//go:build unix

package color

import (
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// queryBackground asks the terminal for its background color with OSC 11.
// A primary device attributes query (DA1) follows it. Every terminal answers
// DA1, so when that answer arrives without an OSC 11 answer before it, the
// terminal does not support the query and there is no need to wait for the
// timeout.
func queryBackground(timeout time.Duration) (Color, bool) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		return None, false
	}
	defer tty.Close() //nolint:errcheck // read-only use, nothing to flush.

	fd := int(tty.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return None, false
	}
	defer term.Restore(fd, state) //nolint:errcheck // nothing to do if restoring fails.

	if _, err := tty.WriteString("\x1b]11;?\x1b\\\x1b[c"); err != nil {
		return None, false
	}

	deadline := time.Now().Add(timeout)

	var resp strings.Builder

	buf := make([]byte, 256)

	for {
		left := time.Until(deadline)
		if left <= 0 {
			return None, false
		}

		fds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}} //nolint:gosec // fd fits in int32.

		n, err := unix.Poll(fds, int(left.Milliseconds())+1)
		if err != nil && err != unix.EINTR {
			return None, false
		}

		if n <= 0 {
			continue
		}

		read, err := unix.Read(fd, buf)
		if err != nil || read <= 0 {
			return None, false
		}

		resp.Write(buf[:read])

		// The DA1 answer ends the exchange: "\x1b[?...c".
		s := resp.String()
		if i := strings.LastIndex(s, "\x1b[?"); i >= 0 && strings.HasSuffix(s, "c") {
			return parseOSC11(s[:i])
		}
	}
}
