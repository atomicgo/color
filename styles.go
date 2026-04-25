package color

// Color shortcut functions apply common ANSI styles to text.
var (
	Black       = NewStyle(ANSIBlack, nil).Sprint
	BrightBlack = NewStyle(ANSIBrightBlack, nil).Sprint

	Red       = NewStyle(ANSIRed, nil).Sprint
	BrightRed = NewStyle(ANSIBrightRed, nil).Sprint

	Green       = NewStyle(ANSIGreen, nil).Sprint
	BrightGreen = NewStyle(ANSIBrightGreen, nil).Sprint

	Yellow       = NewStyle(ANSIYellow, nil).Sprint
	BrightYellow = NewStyle(ANSIBrightYellow, nil).Sprint

	Blue       = NewStyle(ANSIBlue, nil).Sprint
	BrightBlue = NewStyle(ANSIBrightBlue, nil).Sprint

	Magenta       = NewStyle(ANSIMagenta, nil).Sprint
	BrightMagenta = NewStyle(ANSIBrightMagenta, nil).Sprint

	Cyan       = NewStyle(ANSICyan, nil).Sprint
	BrightCyan = NewStyle(ANSIBrightCyan, nil).Sprint

	White       = NewStyle(ANSIWhite, nil).Sprint
	BrightWhite = NewStyle(ANSIBrightWhite, nil).Sprint

	Success = NewStyle(ANSIBrightGreen, nil).Sprint
	Info    = NewStyle(ANSIBrightBlue, nil).Sprint
	Warning = NewStyle(ANSIBrightYellow, nil).Sprint
	Error   = NewStyle(ANSIBrightRed, nil).Sprint
	Fatal   = NewStyle(ANSIBrightRed, nil, Bold).Sprint
)
