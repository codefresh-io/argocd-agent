package logger

import (
	"github.com/gookit/color"
)

func Warning(message string) {
	color.Yellow.Println(message)
	_, _ = color.Reset()
}

func Error(message string) {
	color.Red.Println(message)
	_, _ = color.Reset()
}

func Success(message string) {
	color.Green.Println(message)
	_, _ = color.Reset()
}
