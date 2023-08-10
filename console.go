package log

import (
	"fmt"
)

func wrap(content string, _color Color) string {
	return string(_color) + content + string(RESET)
}

func printc(content string, _color Color) {
	fmt.Print(string(_color) + content + string(RESET))
}

func printcln(content string, _color Color) {
	fmt.Println(string(_color) + content + string(RESET))
}

func printcf(format string, _color Color, params ...interface{}) {
	formatted := fmt.Sprintf(format, params...)
	wrapped := wrap(formatted, _color)
	print(wrapped)
}

func sprintcf(format string, _color Color, params ...any) string {
	formatted := fmt.Sprintf(format, params...)
	wrapped := wrap(formatted, _color)
	return wrapped
}
