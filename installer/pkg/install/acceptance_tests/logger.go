package acceptance_tests

import (
	"github.com/gookit/color"
)

func info(message string){
	color.Println(message)
}

func success(message string){
	green := color.New(color.Green, color.Bold)
	green.Print("√ ")
	_, _ = color.Reset()
	color.Println(message)
}

func failure(message string){
	red := color.New(color.Red, color.Bold)
	red.Print("x ")
	_, _ = color.Reset()
	color.Println(message)
}
