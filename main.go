package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}


var supportedCommands map[string]cliCommand
func main() {
	supportedCommands = map[string]cliCommand{
		"exit":{
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help":{
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for i:=0; ; i++{
		fmt.Print("Pokedex > ")
		scanner.Scan()
		cleanInput := cleanInput(scanner.Text())
		firstWord := cleanInput[0]
		cmd, exists := supportedCommands[firstWord]
		if !exists{
			fmt.Println("Unknown Command")
		}else{
			err := cmd.callback()
			if err != nil{
				fmt.Printf("An error occurred: %v\n", err) 
			}
		}

	

	}
}

func cleanInput(text string) []string{
	formattedString := strings.ToLower(strings.TrimSpace(text))
	splitedString := strings.Fields(formattedString)
	return splitedString

}

func commandExit() error{
	defer os.Exit(0)
	println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp() error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _,cmd := range supportedCommands{
		fmt.Printf("%v: %v\n",cmd.name,cmd.description)
	}
	return nil
}