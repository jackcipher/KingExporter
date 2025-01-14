package display

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Style defines the visual appearance of displayed messages
type Style struct {
	TopChar    string
	BottomChar string
	Width      int
}

// DefaultStyles provides predefined styles for different message types
var DefaultStyles = struct {
	Standard Style
	Error    Style
	Input    Style
	Exit     Style
}{
	Standard: Style{TopChar: "*", BottomChar: "*", Width: 60},
	Error:    Style{TopChar: "!", BottomChar: "!", Width: 60},
	Input:    Style{TopChar: "+-", BottomChar: "+-", Width: 30},
	Exit:     Style{TopChar: "-*", BottomChar: "-*", Width: 30},
}

// Printer handles the output operations
type Printer struct {
	out    io.Writer
	errOut io.Writer
}

// NewPrinter creates a new Printer instance with custom output writers
func NewPrinter(out, errOut io.Writer) *Printer {
	return &Printer{
		out:    out,
		errOut: errOut,
	}
}

// DefaultPrinter is the default instance using standard output and error
var DefaultPrinter = NewPrinter(os.Stdout, os.Stderr)

// getSeparator creates a separator line using the specified character and width
func getSeparator(s string, width int) string {
	if len(s) > 0 {
		repeat := width / len(s)
		if repeat > 0 {
			return strings.Repeat(s, repeat)
		}
	}
	return strings.Repeat("-", width) // fallback
}

// formatMessage formats a message with the given style
func formatMessage(style Style, format string, args ...interface{}) string {
	topSep := getSeparator(style.TopChar, style.Width)
	bottomSep := getSeparator(style.BottomChar, style.Width)
	message := fmt.Sprintf(format, args...)

	var builder strings.Builder
	builder.WriteString(topSep + "\n")
	builder.WriteString(message + "\n")
	builder.WriteString(bottomSep)

	return builder.String()
}

// Print displays a formatted message using the standard style
func (p *Printer) Print(format string, args ...interface{}) {
	fmt.Fprintln(p.out, formatMessage(DefaultStyles.Standard, format, args...))
}

// PrintError displays a formatted error message
func (p *Printer) PrintError(format string, args ...interface{}) {
	fmt.Fprintln(p.errOut, formatMessage(DefaultStyles.Error, format, args...))
}

// PrintInput displays an input prompt
func (p *Printer) PrintInput(format string, args ...interface{}) {
	fmt.Fprintln(p.out, formatMessage(DefaultStyles.Input, format, args...))
	fmt.Fprint(p.out, "> ")
}

// Exit displays a message and exits with the specified code
func (p *Printer) Exit(code int, format string, args ...interface{}) {
	fmt.Fprintln(p.out, formatMessage(DefaultStyles.Exit, format, args...))
	os.Exit(code)
}

// Convenience functions using DefaultPrinter

func Print(format string, args ...interface{}) {
	DefaultPrinter.Print(format, args...)
}

func PrintError(format string, args ...interface{}) {
	DefaultPrinter.PrintError(format, args...)
}

func PrintInput(format string, args ...interface{}) {
	DefaultPrinter.PrintInput(format, args...)
}

func Exit(code int, format string, args ...interface{}) {
	DefaultPrinter.Exit(code, format, args...)
}

// ExitSuccess is a convenience function for successful exit
func ExitSuccess(format string, args ...interface{}) {
	Exit(0, format, args...)
}

// ExitError is a convenience function for error exit
func ExitError(format string, args ...interface{}) {
	Exit(1, format, args...)
}

func FormatBytes(bytes int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
		PB = 1024 * TB
	)

	switch {
	case bytes >= PB:
		return fmt.Sprintf("%.2f PB", float64(bytes)/float64(PB))
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
