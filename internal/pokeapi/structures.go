package pokeapi

type LocationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (resp LocationResponse) GetNames() []string {
	names := make([]string, len(resp.Results))
	for i, location := range resp.Results {
		names[i] = location.Name
	}
	return names
}
