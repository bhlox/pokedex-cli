package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bhlox/pokedex-cli/internal/pokedex"
)

// type Config struct {
// 	previousURL string
// 	nextURL string
// }

type LocationData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationSpecificData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetNewLocations(config *pokedex.Config,params []string) error {    
	cache,ok := config.PokeapiClient.Cache.Get(config.NextURL)
	if ok {
		locationData := &LocationData{}
    	if err := json.Unmarshal(cache, locationData); err != nil {
        	return fmt.Errorf("failed to unmarshal JSON: %w", err)
    	}

    	for _, result := range locationData.Results {
        	fmt.Println(result.Name)
    	}
    	config.NextURL = locationData.Next
    	config.PreviousURL = locationData.Previous
    	return nil
	}

    res, err := http.Get(config.NextURL)
    if err != nil {
        return fmt.Errorf("failed to get new locations: %w", err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }
    if res.StatusCode > 299 {
        return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
    }
	config.PokeapiClient.Cache.Add(config.NextURL,body)
    locationData := &LocationData{}
    if err := json.Unmarshal(body, locationData); err != nil {
        return fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    for _, result := range locationData.Results {
        fmt.Println(result.Name)
    }
    config.NextURL = locationData.Next
    config.PreviousURL = locationData.Previous
    return nil
}

func GetPrevLocations(config *pokedex.Config,params []string) error {
    apiLink := config.PreviousURL
    if apiLink == "" {
        return errors.New("there are no previous locations to display")
    }
	cache,ok := config.PokeapiClient.Cache.Get(apiLink)
	if ok {
		locationData := &LocationData{}
    	if err := json.Unmarshal(cache, locationData); err != nil {
        	return fmt.Errorf("failed to unmarshal JSON: %w", err)
    	}

    	for _, result := range locationData.Results {
        	fmt.Println(result.Name)
    	}
    	config.NextURL = locationData.Next
    	config.PreviousURL = locationData.Previous
    	return nil
	}

    res, err := http.Get(apiLink)
    if err != nil {
        return fmt.Errorf("failed to get previous locations: %w", err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }
    if res.StatusCode > 299 {
        return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
    }
	config.PokeapiClient.Cache.Add(apiLink,body)
    locationData := &LocationData{}
    if err := json.Unmarshal(body, locationData); err != nil {
        return fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    for _, result := range locationData.Results {
        fmt.Println(result.Name)
    }
    config.NextURL = locationData.Next
    config.PreviousURL = locationData.Previous
    return nil
}
