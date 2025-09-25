package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pepegakac123/pokedexcli/internal/pokecache"
)

type LocationResponse struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []LocationResult `json:"results"`
}

type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var globalCache *pokecache.Cache

func init() {
	globalCache = pokecache.NewCache(7 * time.Second)
}

func GetLocationAreas(url string) (*LocationResponse, error) {
	cachedValue, ok := globalCache.Get(url)
	var data []byte
	if ok {
		data = cachedValue
	} else {
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return nil, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, data)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		globalCache.Add(url, data)

	}
	var locationResponse LocationResponse
	err := json.Unmarshal(data, &locationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(locationResponse.Results) <= 0 {
		return nil, fmt.Errorf("no locations in the response")
	}

	return &locationResponse, nil
}
