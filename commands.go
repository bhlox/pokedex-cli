package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bhlox/pokedex-cli/internal/pokedex"
	"github.com/bhlox/pokedex-cli/internal/pokedex/location"
	"github.com/bhlox/pokedex-cli/internal/pokedex/pokemon"
	"github.com/bhlox/pokedex-cli/internal/pokedex/utils"
)

func getCommands()map[string]cliCommands{
	return map[string]cliCommands{
		"help": {name: "help", description: "Displays the list of commands", callback: CommandHelp},
		"exit": {name: "exit", description: "exit the pokedex", callback: CommandExit},
		"map":{name:"map",description: "display 20 locations",callback: location.GetNewLocations},
		"mapb":{name:"mapb",description: "display previous 20 locations",callback: location.GetPrevLocations},
		"explore":{name:"explore",description: "explore a location area",callback: CommandExplore},
		"catch":{name:"catch",description: "catch a pokemon of your choice",callback: CommandCatch},
		"pokedex":{name: "pokedex",description: "Prints a list of caught pokemon",callback: LogCaughtPokemon},
		"inspect":{name: "inspect",description: "checkts the stats of caught pokemon",callback: InspectPokemon },
	}
}

func CommandExit(c *pokedex.Config,params []string) error {
	fmt.Println("Exiting the pokedex")
	os.Exit(0)
	return nil
}

func CommandHelp (c *pokedex.Config,params []string)error {
	fmt.Println()
	fmt.Println("displaying list of commands")
	fmt.Println()
	for k, v := range getCommands() {
		fmt.Printf("%v: %v\n", k, v.description)
	}
	fmt.Println()
	return nil
}

func CommandExplore (config *pokedex.Config,params []string)error {
	fmt.Println()
	if len(params) == 0 {
		return errors.New("no location was provided for you to explore")
	}
	locationQuery := params[0]
	apiLink := pokedex.LocationBaseUrl + locationQuery
	data,ok := config.PokeapiClient.Cache.Get(apiLink)
	if ok {
		locationData := &location.LocationSpecificData{}
    	if err := json.Unmarshal(data, locationData); err != nil {
        	return fmt.Errorf("failed to unmarshal JSON: %w", err)
	    }
		fmt.Printf("exploring area of %v \n",locationData.Location.Name)
		fmt.Println("Found pokemon:")
		for _,pokemon := range locationData.PokemonEncounters {
			fmt.Printf("- %v\n",pokemon.Pokemon.Name)
		}
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
    if res.StatusCode > 400 {
		fmt.Printf("Failed to get details of location: %v.\n",locationQuery)
        return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
    }
	config.PokeapiClient.Cache.Add(apiLink,body)
    locationData := &location.LocationSpecificData{}
    if err := json.Unmarshal(body, locationData); err != nil {
        return fmt.Errorf("failed to unmarshal JSON: %w", err)
    }
	fmt.Printf("exploring area of %v \n",locationData.Location.Name)
	utils.DummyLoading()
	fmt.Println("Found pokemon:")
	for _,pokemon := range locationData.PokemonEncounters {
		fmt.Printf("- %v\n",pokemon.Pokemon.Name)
	}
	return nil
}

func CommandCatch(config *pokedex.Config,params []string) error {
	fmt.Println()
	if len(params) == 0 {
		return errors.New("no pokemon was provided to catch")
	}
	pokemonQuery := params[0]
	apiLink := "https://pokeapi.co/api/v2/pokemon/" + pokemonQuery
	res, err := http.Get(apiLink)
	if err != nil {
        return fmt.Errorf("failed to get previous locations: %w", err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }
    if res.StatusCode > 400 {
		fmt.Printf("Failed to get details of pokemon: %v.\n",pokemonQuery)
        return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
    }
	details := &pokemon.Details{}
	if err := json.Unmarshal(body, details); err != nil {
        return fmt.Errorf("failed to unmarshal JSON: %w", err)
    }
	// #TODO add some additional factors. e.g. rateToCatchPokemon becomes dependent on the pokemons base experience. if chance of success is low, prompt the user if wants to proceed or not. if proceed (has a chance to increase success rate by 10-50%)
	fmt.Printf("Attemping to catch %v\n",details.Name)
	utils.DummyLoading()
	rateToCatchPokemon := utils.GenerateRandomNum(15,100)
	trainerSuccessRate := utils.GenerateRandomNum(0,100)
	fmt.Printf("success rate to catch %v: %v, trainer's success rate: %v\n",details.Name,rateToCatchPokemon,trainerSuccessRate)
	if trainerSuccessRate + rateToCatchPokemon >= 100 {
		config.Pokedex[details.Name] = *details
		fmt.Printf("Successfully caught %v!\n",details.Name)
	} else {
		fmt.Printf("%v managed to break out and has fled\n",details.Name)
	}
	return nil
}

func LogCaughtPokemon(config *pokedex.Config, params []string) error {
	println("Your caught pokemon:")
	pokedex := config.Pokedex
	for _,v := range pokedex {
		fmt.Printf("- %v\n",v.Name)
	}
	return nil
}

func InspectPokemon(config *pokedex.Config, params []string) error {
	pokemonQueried := params[0]
	if pokemon, ok := config.Pokedex[pokemonQueried]; ok {
		fmt.Printf("Name: %v\n",pokemon.Name)
		fmt.Printf("Height: %v\n",pokemon.Height)
		fmt.Printf("Weight: %v \n",pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _,v := range pokemon.Stats {
			fmt.Printf("	%v: %v\n",v.Stat.Name,v.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _,v := range pokemon.Types {
			fmt.Printf("	-%v\n",v.Type.Name)
		}
		return nil
	}
	
	return fmt.Errorf("you haven't caught %v yet",pokemonQueried)
}