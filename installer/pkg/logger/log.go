package logger

import "fmt"

var (
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
)

func Warning(message string) {
	fmt.Println(colorYellow, message)
}

func Error(message string) {
	fmt.Println(colorRed, message)
}

func Success(message string) {
	fmt.Println(colorGreen, message)
}
