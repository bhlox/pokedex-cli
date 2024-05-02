package pokedex

import (
	"net/http"
	"time"

	"github.com/bhlox/pokedex-cli/internal/pokedex/cache"
)

type Client struct {
	Cache      cache.Cache
	HttpClient http.Client
}

// #TODO remove htppclient?
func NewClient(cacheInterval time.Duration)Client{
	return Client{Cache: *cache.NewCache(cacheInterval),HttpClient: *http.DefaultClient}
}