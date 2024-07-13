package color

import (
	"golang.org/x/sys/windows"
)

func init() {
	handle, _ := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)

	var mode uint32
	windows.GetConsoleMode(handle, &mode)

	// Reference: https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
		vtpmode := mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		windows.SetConsoleMode(handle, vtpmode)
	}
}
