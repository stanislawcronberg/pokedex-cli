package pokeapi

import (
	"github.com/stanislawcronberg/pokedex-cli/internal/pokecache"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

type SessionState struct {
	Client   Client
	Cache    *pokecache.Cache
	Next     *string
	Previous *string
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
	}
}
