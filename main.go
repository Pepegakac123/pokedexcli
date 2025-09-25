package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string{
	formattedString := strings.ToLower(strings.TrimSpace(text))
	splitedString := strings.Fields(formattedString)
	return splitedString

}