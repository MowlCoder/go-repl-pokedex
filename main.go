package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/MowlCoder/go-repl-pokedex/pokeapi"
	"github.com/MowlCoder/go-repl-pokedex/pokestore"
	"math/rand"
	"os"
	"strings"
)

type CliCommand struct {
	name        string
	description string
	callback    func(args []string, config *Config) error
}

type Config struct {
	NextMapUrl string
	PrevMapUrl string
}

type commands map[string]CliCommand

func (c commands) addCommand(command CliCommand) {
	c[command.name] = command
}

var pokemonStore = pokestore.NewPokemonStore()

func onMapCommand(args []string, config *Config) error {
	if config.NextMapUrl == "" {
		return errors.New("you are already on a last page")
	}

	locationResponse, err := pokeapi.GetLocations(config.NextMapUrl)

	if err != nil {
		return err
	}

	config.PrevMapUrl = locationResponse.Previous
	config.NextMapUrl = locationResponse.Next

	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func onMapBCommand(args []string, config *Config) error {
	if config.PrevMapUrl == "" {
		return errors.New("you are already on a first page")
	}

	locationResponse, err := pokeapi.GetLocations(config.PrevMapUrl)

	if err != nil {
		return err
	}

	config.PrevMapUrl = locationResponse.Previous
	config.NextMapUrl = locationResponse.Next

	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func onExploreCommand(args []string, config *Config) error {
	if len(args) != 1 {
		return errors.New("you have to provide name of exploring area")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	locationAreaResponse, err := pokeapi.GetLocationArea(areaName)

	if err != nil {
		return err
	}

	fmt.Println("Found pokemon:")

	for _, pokemon := range locationAreaResponse.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func onCatchCommand(args []string, config *Config) error {
	if len(args) != 1 {
		return errors.New("you have to provide name of pokemon")
	}

	pokemonName := args[0]
	fmt.Printf("Catching %s...\n", pokemonName)

	pokemonInfo, err := pokeapi.GetPokemonInfo(pokemonName)

	if err != nil {
		return err
	}

	randInt := rand.Intn(pokemonInfo.BaseExperience)

	if pokemonInfo.BaseExperience-randInt <= 30 {
		types := make([]string, len(pokemonInfo.Types))

		for idx, pokeType := range pokemonInfo.Types {
			types[idx] = pokeType.Type.Name
		}

		stats := make(map[string]int)

		for _, stat := range pokemonInfo.Stats {
			stats[stat.Stat.Name] = stat.BaseStat
		}

		pokemonStore.Add(pokestore.Pokemon{
			ID:     pokemonInfo.ID,
			Name:   pokemonInfo.Name,
			Height: pokemonInfo.Height,
			Weight: pokemonInfo.Weight,
			Stats:  stats,
			Types:  types,
		})

		fmt.Printf("you catched %s\n", pokemonName)
	} else {
		fmt.Printf("%s is escaped\n", pokemonName)
	}

	return nil
}

func onInspectCommand(args []string, config *Config) error {
	if len(args) != 1 {
		return errors.New("you have to provide name of pokemon")
	}

	pokemonName := args[0]
	pokemon, ok := pokemonStore.Get(pokemonName)

	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	fmt.Println(pokemon)

	return nil
}

func onPokeDexCommand(args []string, config *Config) error {
	pokemonNames := pokemonStore.GetAllKeys()

	if len(pokemonNames) == 0 {
		fmt.Println("your pokedex is empty")
		return nil
	}

	fmt.Println("Your Pokedex:")

	for _, name := range pokemonNames {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}

func onHelpCommand(args []string, config *Config) error {
	fmt.Println("================PokeDex Cli================")

	for key, val := range replCommands {
		fmt.Printf("%s: %s\n", key, val.description)
	}

	fmt.Println("===========================================")

	return nil
}

func onExitCommand(args []string, config *Config) error {
	os.Exit(0)
	return nil
}

var replCommands = commands{}
var replConfig = Config{
	NextMapUrl: "https://pokeapi.co/api/v2/location-area",
}

func initReplCommands() {
	replCommands.addCommand(CliCommand{
		name:        "map",
		description: "display a locations (forward)",
		callback:    onMapCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "mapb",
		description: "display a locations (back)",
		callback:    onMapBCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "explore",
		description: "explore location for pokemons",
		callback:    onExploreCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "catch",
		description: "try to catch pokemon",
		callback:    onCatchCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "inspect",
		description: "display information about pokemon",
		callback:    onInspectCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "pokedex",
		description: "display list of your pokemons",
		callback:    onPokeDexCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "help",
		description: "list a command and their descriptions",
		callback:    onHelpCommand,
	})

	replCommands.addCommand(CliCommand{
		name:        "exit",
		description: "exit from pokedex cli",
		callback:    onExitCommand,
	})
}

func main() {
	initReplCommands()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()

		userInput := scanner.Text()
		fields := strings.Fields(userInput)

		if command, ok := replCommands[fields[0]]; ok {
			err := command.callback(fields[1:], &replConfig)

			if err != nil {
				fmt.Printf("error: %v\n", err)
			}

			continue
		}

		fmt.Println("Unknown command. Type help to get list of a commands")
	}
}
