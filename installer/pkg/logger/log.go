package logger

import (
	"github.com/gookit/color"
	"strconv"
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

func Info(message string) {
	color.Println(message)
}

func SuccessTest(message string) {
	green := color.New(color.Green, color.Bold)
	green.Print("√ ")
	_, _ = color.Reset()
	color.Println(message)
}

func FailureTest(message string) {
	red := color.New(color.Red, color.Bold)
	red.Print("x ")
	_, _ = color.Reset()
	color.Println(message)
}

func Summary(number int, message string, value string) {
	color.Print("    " + strconv.Itoa(number) + ". " + message + ": ")
	color.Cyan.Println(value)
	_, _ = color.Reset()
}
