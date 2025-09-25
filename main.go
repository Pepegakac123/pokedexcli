package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Pepegakac123/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	prevUrl string
	nextUrl string
}

type LocationResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []LocationResults `json:"results"`
}

type LocationResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var config Config
var supportedCommands map[string]cliCommand

func main() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays the names of 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the names of 20 previous location areas in the Pokemon world.",
			callback:    commandMapb,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		cleanInput := cleanInput(scanner.Text())
		if len(cleanInput) > 0 {
			firstWord := cleanInput[0]
			cmd, exists := supportedCommands[firstWord]
			if !exists {
				fmt.Println("Unknown Command")
			} else {
				err := cmd.callback(&config)
				if err != nil {
					fmt.Printf("An error occurred: %v\n", err)
				}
			}
		}

	}
}

func cleanInput(text string) []string {
	formattedString := strings.ToLower(strings.TrimSpace(text))
	splitedString := strings.Fields(formattedString)
	return splitedString

}

func commandExit(cfg *Config) error {
	defer os.Exit(0)
	println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range supportedCommands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
func commandMap(cfg *Config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.nextUrl != "" {
		url = cfg.nextUrl
	}

	response, err := pokeapi.GetLocationAreas(url)
	if err != nil {
		return err
	}

	updateConfigUrls(cfg, *response)

	for _, location := range response.Results {
		fmt.Printf("%v\n", location.Name)
	}

	return nil
}

func commandMapb(cfg *Config) error {
	if cfg.prevUrl == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	response, err := pokeapi.GetLocationAreas(cfg.prevUrl)
	if err != nil {
		return err
	}

	updateConfigUrls(cfg, *response)

	for _, location := range response.Results {
		fmt.Printf("%v\n", location.Name)
	}

	return nil
}

func updateConfigUrls(cfg *Config, response pokeapi.LocationResponse) {
	if response.Next != nil {
		cfg.nextUrl = *response.Next
	} else {
		cfg.nextUrl = ""
	}
	if response.Previous != nil {
		cfg.prevUrl = *response.Previous
	} else {
		cfg.prevUrl = ""
	}
}
