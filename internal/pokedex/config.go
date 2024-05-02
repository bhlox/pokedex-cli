package pokedex

import (
	"time"

	"github.com/bhlox/pokedex-cli/internal/pokedex/pokemon"
)

const LocationBaseUrl = "https://pokeapi.co/api/v2/location-area/"

type Config struct {
	PokeapiClient Client
	PreviousURL   string
	NextURL       string
	Pokedex       map[string]pokemon.Details
}

func InitConfig() *Config{
	client := NewClient(5 * time.Minute)
	return &Config{PokeapiClient: client,NextURL: LocationBaseUrl,Pokedex: make(map[string]pokemon.Details)}
}
