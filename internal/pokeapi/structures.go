package pokeapi

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// LocationAreaResponse stores part of the information from the
// location-area endpoint of the PokeAPI, we are only interested
// in the Pokémon we can find in each location
type LocationAreaResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	}
}

func (resp LocationAreasResponse) GetLocationNames() []string {
	names := make([]string, len(resp.Results))
	for i, location := range resp.Results {
		names[i] = location.Name
	}
	return names
}

func (resp LocationAreaResponse) GetPokemonNames() []string {
	names := make([]string, len(resp.PokemonEncounters))
	for i, encounter := range resp.PokemonEncounters {
		names[i] = encounter.Pokemon.Name
	}
	return names
}
