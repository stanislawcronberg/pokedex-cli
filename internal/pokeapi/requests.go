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

func unmarshalLocationResponse(data []byte) (LocationAreasResponse, error) {
	var locations LocationAreasResponse
	err := json.Unmarshal(data, &locations)
	if err != nil {
		return LocationAreasResponse{}, fmt.Errorf("error unmarshalling data: %s", err)
	}
	return locations, nil
}

func (c *Client) GetLocationAreasResponse(newURL *string, cache *pokecache.Cache) (LocationAreasResponse, error) {

	fullURL := buildURL(newURL, "/location")

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

	locations, err := unmarshalLocationResponse(data)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locations, nil
}
