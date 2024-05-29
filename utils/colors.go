package utils

import "fmt"

const (
	Reset    = "\033[0m"
	Red      = "\033[31m"
	Green    = "\033[32m"
	Yellow   = "\033[33m"
	Blue     = "\033[34m"
	Purple   = "\033[35m"
	Cyan     = "\033[36m"
	White    = "\033[37m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgCyan   = "\033[46m"
	BgWhite  = "\033[47m"
)

func Colorize(color string, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

func RedText(text string) string {
	return Colorize(Red, text)
}

func GreenText(text string) string {
	return Colorize(Green, text)
}

func YellowText(text string) string {
	return Colorize(Yellow, text)
}

func BlueText(text string) string {
	return Colorize(Blue, text)
}

func PurpleText(text string) string {
	return Colorize(Purple, text)
}

func CyanText(text string) string {
	return Colorize(Cyan, text)
}

func WhiteText(text string) string {
	return Colorize(White, text)
}

func BgRedText(text string) string {
	return Colorize(BgRed, text)
}

func BgGreenText(text string) string {
	return Colorize(BgGreen, text)
}

func BgYellowText(text string) string {
	return Colorize(BgYellow, text)
}

func BgBlueText(text string) string {
	return Colorize(BgBlue, text)
}

func BgPurpleText(text string) string {
	return Colorize(BgPurple, text)
}

func BgCyanText(text string) string {
	return Colorize(BgCyan, text)
}

func BgWhiteText(text string) string {
	return Colorize(BgWhite, text)
}

func CustomColorText(fgColor string, bgColor string, text string) string {
	return fmt.Sprintf("%s%s%s%s", fgColor, bgColor, text, Reset)
}
