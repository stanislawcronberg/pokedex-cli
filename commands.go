package main

import (
	"encoding/json"
	"fmt"
	"github.com/stanislawcronberg/pokedex-cli/internal/pokeapi"
	"os"
)

func helpCallback(conf *pokeapi.Config, params []string) error {
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

func exitCallback(conf *pokeapi.Config, params []string) error {
	defer os.Exit(0)
	return nil
}

func printItems(items []string) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func updateConfig(conf *pokeapi.Config, locationResponse *pokeapi.LocationAreasResponse) {
	conf.Next = locationResponse.Next
	conf.Previous = locationResponse.Previous
}

func getLocationResponse(conf *pokeapi.Config, url *string) (pokeapi.LocationAreasResponse, error) {
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

func nextLocationsCallback(conf *pokeapi.Config, params []string) error {
	locationResponse, err := getLocationResponse(conf, conf.Next)
	if err != nil {
		return err
	}

	updateConfig(conf, &locationResponse)
	printItems(locationResponse.GetLocationNames())
	return nil
}

func previousLocationsCallback(conf *pokeapi.Config, params []string) error {
	locationResponse, err := getLocationResponse(conf, conf.Previous)
	if err != nil {
		return err
	}

	updateConfig(conf, &locationResponse)
	printItems(locationResponse.GetLocationNames())
	return nil
}

func locationAreaPokemonsCallback(conf *pokeapi.Config, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("no location provided")
	}

	locationName := params[0]

	locationResponse, err := conf.Client.GetLocationAreaResponse(locationName, conf.Cache)
	if err != nil {
		return err
	}

	printItems(locationResponse.GetPokemonNames())

	return nil
}
