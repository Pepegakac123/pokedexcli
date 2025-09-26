package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/Pepegakac123/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

type Pokedex struct {
	pokedex map[string]pokeapi.Pokemon
}

type Config struct {
	prevUrl string
	nextUrl string
}

var pokedex map[string]pokeapi.Pokemon
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
		"explore": {
			name:        "explore",
			description: "It takes the name of a location area as an argument and lists the names of the Pokemon that can be found there.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "It take a name of the pokemon and attempts to catch it",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "It take a name of the pokemon and show his stats only if it is in your pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "It lists all of your pokemon",
			callback:    commandPokedex,
		},
	}
	pokedex = make(map[string]pokeapi.Pokemon)
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
				parametr := ""
				if len(cleanInput) > 1 {
					parametr = cleanInput[1]
				}
				err := cmd.callback(&config, parametr)
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

func commandExit(cfg *Config, location string) error {
	defer os.Exit(0)
	println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(cfg *Config, location string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range supportedCommands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
func commandMap(cfg *Config, location string) error {
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

func commandMapb(cfg *Config, location string) error {
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

func commandExplore(cfg *Config, location string) error {
	if location == "" {
		return fmt.Errorf("The command must have <location> parametr")
	}
	response, err := pokeapi.GetLocationEncounters(location)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %v...\n", location)
	fmt.Println("Found Pokemon:")

	for _, pokemon := range response.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}
	return nil
}
func commandCatch(cfg *Config, pokemon string) error {
	if pokemon == "" {
		return fmt.Errorf("The command must have <pokemon> parametr")
	}
	response, err := pokeapi.GetPokemon(pokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)
	szansa_zlapania := int(math.Max(5, 100-float64(response.BaseExperience)/3))
	random := rand.Intn(100)
	if random < szansa_zlapania {
		fmt.Printf("%v was caught!\n", pokemon)
		addPokemonToPokedex(response, pokemon)

	} else {
		fmt.Printf("%v escaped!\n", pokemon)
	}
	return nil
}

func addPokemonToPokedex(pokemon *pokeapi.Pokemon, pokemonName string) {
	if _, ok := pokedex[pokemonName]; !ok {
		pokedex[pokemonName] = *pokemon
	}
}

func commandInspect(cfg *Config, pokemon string) error {
	if pokemon == "" {
		return fmt.Errorf("The command must have <pokemon> parametr")
	}

	if _, ok := pokedex[pokemon]; !ok {
		return fmt.Errorf("You have not caught that pokemon")
	}

	printPokemon(pokemon)
	return nil
}

func printPokemon(pokemon string) {
	pokemonInfo := pokedex[pokemon]
	fmt.Printf("Name: %v\n", pokemonInfo.Name)
	fmt.Printf("Height: %v\n", pokemonInfo.Height)
	fmt.Printf("Weight: %v\n", pokemonInfo.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemonInfo.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemonInfo.Types {
		fmt.Printf("  - %v\n", pokemonType.Type.Name)
	}
}

func commandPokedex(cfg *Config, parametr string) error {
	fmt.Println("Your Pokedex:")
	for pokemon, _ := range pokedex {
		fmt.Printf(" - %v\n", pokemon)
	}
	return nil
}
