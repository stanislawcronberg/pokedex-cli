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

func printLocations(locations []string) {
	for _, location := range locations {
		fmt.Println(location)
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

	locationResponse, err := conf.Client.ListLocations(url, conf.Cache)
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
	printLocations(locationResponse.GetNames())
	return nil
}

func previousLocationsCallback(conf *pokeapi.Config, params []string) error {
	locationResponse, err := getLocationResponse(conf, conf.Previous)
	if err != nil {
		return err
	}

	updateConfig(conf, &locationResponse)
	printLocations(locationResponse.GetNames())
	return nil
}
