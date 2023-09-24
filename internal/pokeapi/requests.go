package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/stanislawcronberg/pokedex-cli/internal/pokecache"
	"io"
	"log"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

func buildURL(newURL *string, endpoint string) string {
	fullURL := baseURL + endpoint
	if newURL != nil {
		fullURL = *newURL
	}
	return fullURL
}

func executeGetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %s", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status code %d received from API", resp.StatusCode)
	}

	return resp, nil
}

func processResponse(resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	return data, nil
}

func unmarshalLocationAreasResponse(data []byte) (LocationAreasResponse, error) {
	var locations LocationAreasResponse
	err := json.Unmarshal(data, &locations)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("error unmarshalling data: %s", err)
	}
	return locations, nil
}

func unmarshalLocationAreaResponse(data []byte) (LocationAreaResponse, error) {
	var location LocationAreaResponse
	err := json.Unmarshal(data, &location)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error unmarshalling data: %s", err)
	}
	return location, nil
}

func unmarshalPokemonResponse(data []byte) (PokemonResponse, error) {
	var pokemon PokemonResponse
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return PokemonResponse{}, fmt.Errorf("error unmarshalling data: %s", err)
	}
	return pokemon, nil
}

func (c *Client) GetLocationAreasResponse(newURL *string, cache *pokecache.Cache) (LocationAreasResponse, error) {

	fullURL := buildURL(newURL, "/location-area")

	response, err := executeGetRequest(fullURL)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	data, err := processResponse(response)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	err = cache.Add(&fullURL, data)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	locations, err := unmarshalLocationAreasResponse(data)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locations, nil
}

func (c *Client) GetLocationAreaResponse(locationName string, cache *pokecache.Cache) (LocationAreaResponse, error) {

	fullURL := buildURL(nil, "/location-area/"+locationName)

	response, err := executeGetRequest(fullURL)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	data, err := processResponse(response)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	err = cache.Add(&fullURL, data)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	location, err := unmarshalLocationAreaResponse(data)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return location, nil
}

func (c *Client) GetPokemonResponse(pokemonName string, cache *pokecache.Cache) (PokemonResponse, error) {

	fullURL := buildURL(nil, "/pokemon/"+pokemonName)

	response, err := executeGetRequest(fullURL)
	if err != nil {
		return PokemonResponse{}, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	data, err := processResponse(response)
	if err != nil {
		return PokemonResponse{}, err
	}

	err = cache.Add(&fullURL, data)
	if err != nil {
		return PokemonResponse{}, err
	}

	pokemon, err := unmarshalPokemonResponse(data)
	if err != nil {
		return PokemonResponse{}, err
	}

	return pokemon, nil
}
