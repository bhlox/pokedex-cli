package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bhlox/pokedex-cli/internal/pokedex"
)

type cliCommands struct {
	name        string
	description string
	callback    func(*pokedex.Config, []string) error
}

func Reader(config *pokedex.Config) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the pokedex. Choose from one of the commands to get started")
	fmt.Println()
	commands := getCommands()
	for k, v := range commands{ 
		fmt.Printf("%v: %v\n", k, v.description)
	}
	fmt.Println()
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			return fmt.Errorf("no input provided")
		}
		input := strings.Fields(cleanInput(scanner.Text()))
		commandInput,parameters := input[0],input[1:]
		command, ok := commands[commandInput]
		if !ok {
			fmt.Println()
			fmt.Printf("couldn't find this command: %s\n", input)
			fmt.Println()
			continue
		}
		fmt.Println()
		err := command.callback(config,parameters)
		if err != nil {
			log.Println("Error occurred:", err)
		}
	}
}

func cleanInput(input string)string{
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	return input
}