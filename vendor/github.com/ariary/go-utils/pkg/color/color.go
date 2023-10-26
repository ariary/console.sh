package color

import "fmt"

var (
	Info = Teal
	Warn = Yellow
	Evil = Red
	Good = Green
	Code = Cyan
)

var (
	// Background
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Blue    = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
	Cyan    = Color("\033[1;96m%s\033[0m")

	// Foreground
	RedForeground     = Color("\033[1;41m%s\033[0m")
	GreenForeground   = Color("\033[1;42m%s\033[0m")
	YellowForeground  = Color("\033[1;43m%s\033[0m")
	BlueForeground    = Color("\033[1;44m%s\033[0m")
	MagentaForeground = Color("\033[1;45m%s\033[0m")
	TealForeground    = Color("\033[1;46m%s\033[0m")
	WhiteForeground   = Color("\033[1;47m%s\033[0m")

	// Style
	Bold       = Color("\033[1m%s\033[0m")
	Dim        = Color("\033[2m%s\033[0m")
	Italic     = Color("\033[3m%s\033[0m")
	Underlined = Color("\033[4m%s\033[24m")
	Blink      = Color("\033[5m%s\033[24m")
	Reverse    = Color("\033[7m%s\033[24m")
	Hidden     = Color("\033[8m%s\033[24m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
