package main

import (
	"fmt"

	"github.com/bhlox/pokedex-cli/internal/pokedex"
)


func main() {
	config := pokedex.InitConfig()
	err := Reader(config)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
