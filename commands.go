package main

import (
	"encoding/json"
	"fmt"
	"github.com/stanislawcronberg/pokedex-cli/internal/pokeapi"
	"math/rand"
	"os"
)

func helpCallback(conf *pokeapi.SessionState, args ...string) error {
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("- Welcome to the Pokedex, a CLI tool for looking up Pokemon! -")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("=          Here are the commands you can use                 =")
	fmt.Println("--------------------------------------------------------------")

	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%s - %s\n", command.name, command.description)
	}
	return nil
}

func exitCallback(conf *pokeapi.SessionState, args ...string) error {
	defer os.Exit(0)
	return nil
}

func printItems(items []string) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func updateConfig(conf *pokeapi.SessionState, locationResponse *pokeapi.LocationAreasResponse) {
	conf.Next = locationResponse.Next
	conf.Previous = locationResponse.Previous
}

func getLocationResponse(conf *pokeapi.SessionState, url *string) (pokeapi.LocationAreasResponse, error) {
	var locationResponse pokeapi.LocationAreasResponse

	data, found := conf.Cache.Get(url)
	if found {
		err := json.Unmarshal(data, &locationResponse)
		if err != nil {
			return locationResponse, fmt.Errorf("error unmarshalling data: %s", err)
		}
		return locationResponse, nil
	}

	locationResponse, err := conf.Client.GetLocationAreasResponse(url, conf.Cache)
	if err != nil {
		return locationResponse, fmt.Errorf("error getting locations: %s", err)
	}

	return locationResponse, nil
}

func nextLocationsCallback(conf *pokeapi.SessionState, args ...string) error {
	locationResponse, err := getLocationResponse(conf, conf.Next)
	if err != nil {
		return err
	}

	updateConfig(conf, &locationResponse)
	printItems(locationResponse.GetLocationNames())
	return nil
}

func previousLocationsCallback(conf *pokeapi.SessionState, args ...string) error {
	locationResponse, err := getLocationResponse(conf, conf.Previous)
	if err != nil {
		return err
	}

	updateConfig(conf, &locationResponse)
	printItems(locationResponse.GetLocationNames())
	return nil
}

func locationAreaPokemonsCallback(conf *pokeapi.SessionState, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no location provided")
	}

	locationName := args[0]

	locationResponse, err := conf.Client.GetLocationAreaResponse(locationName, conf.Cache)
	if err != nil {
		return err
	}

	printItems(locationResponse.GetPokemonNames())

	return nil
}

func catchPokemonCallback(conf *pokeapi.SessionState, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon provided")
	}

	pokemonName := args[0]

	pokemonResponse, err := conf.Client.GetPokemonResponse(pokemonName, conf.Cache)
	if err != nil {
		return err
	}

	// Catch the pokemon with a random chance depending on the BaseExperience of the pokemon
	// The higher the BaseExperience, the lower the chance of catching the pokemon
	if rand.Float32() < float32(pokemonResponse.BaseExperience)/500 {
		fmt.Printf("You caught %s!\n", pokemonName)
	} else {
		fmt.Printf("You failed to catch %s!\n", pokemonName)
	}

	return nil
}
