package pkg

import (
	"fmt"
	"sync"
)

var (
	wg = sync.WaitGroup{}
)

func printBoxedMessage(message string) {
	boxWidth := len(message) + 4
	horizontalLine := "+" + repeatChar("-", boxWidth-2) + "+"

	fmt.Println(horizontalLine)
	fmt.Printf("| %s |\n", centerText(message, boxWidth-4))
	fmt.Println(horizontalLine)
}

func repeatChar(char string, times int) string {
	result := ""
	for i := 0; i < times; i++ {
		result += char
	}
	return result
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}

	padding := width - len(text)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return repeatChar(" ", leftPadding) + text + repeatChar(" ", rightPadding)
}
