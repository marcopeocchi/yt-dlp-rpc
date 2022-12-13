package cli

import "fmt"

const (
	// FG
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Reset   = "\033[0m"
	// BG
	BgRed   = "\033[1;41m"
	BgBlue  = "\033[1;44m"
	BgGreen = "\033[1;42m"
)

// Formats a message with the specified ascii escape code, then reset.
func Format(message string, code string) string {
	return fmt.Sprintf("%s%s%s", code, message, Reset)
}
