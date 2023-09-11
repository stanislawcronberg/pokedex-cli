package main

import (
	"bufio"
	"fmt"
	"github.com/stanislawcronberg/pokedex-cli/internal/pokeapi"
	"github.com/stanislawcronberg/pokedex-cli/internal/pokecache"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    exitCallback,
		},
		"help": {
			name:        "help",
			description: "Show a help message",
			callback:    helpCallback,
		},
		"map": {
			name:        "map",
			description: "Lists next locations in the Pokemon world",
			callback:    nextLocationsCallback,
		},
		"mapback": {
			name:        "mapback",
			description: "Lists previous locations in the Pokemon world",
			callback:    previousLocationsCallback,
		},
		"explore": {
			name:        "explore",
			description: "Lists Pokemon in a given location",
			callback:    locationAreaPokemonsCallback,
		},
	}
}

func cleanInput(input string) []string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	words := strings.Split(input, " ")
	return words
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	config := pokeapi.Config{
		Client: pokeapi.NewClient(),
		Cache:  pokecache.NewCache(time.Minute * 5),
	}

	for {
		fmt.Print("pokedex> ")
		scanner.Scan()

		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		command, ok := commands[input[0]]
		if !ok {
			fmt.Printf("pokedex> Unknown command: %s\n", input)
			continue
		}

		err := command.callback(&config, input[1:])
		if err != nil {
			fmt.Printf("pokedex> Error executing command: %s\n", err)
		}
	}
}
