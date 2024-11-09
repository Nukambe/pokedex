package commands

import (
	"errors"
	"fmt"
	"github.com/Nukambe/pokedex/pokeapi"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

var commands = createCommandMap()

func ExecuteCommand(text string) int {
	inputs := strings.Split(text, " ")
	cmd := inputs[0]
	var args []string
	if len(inputs) > 1 {
		args = inputs[1:]
	}

	command, ok := commands[cmd]
	if !ok {
		fmt.Println("invalid command:", text)
		return 2
	}
	if err := command.callback(args); err != nil {
		fmt.Println(err)
		return 1
	}
	if text == "quit" {
		return -1
	}
	return 0
}

func createCommandMap() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"quit": {
			name:        "quit",
			description: "Quit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call to map displays the next 20 locations, and so on. 'map' lets you explore the world of Pokemon.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Similar to the map command, however, instead of displaying the next 20 locations, it displays the previous 20 locations. It's a way to go back.",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "use 'explore [location name or id]' to see all the pokemon available at that location.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "use 'catch [pokemon name or id]' to attempt to catch that pokemon.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "use 'inspect [pokemon]' to see a pokemon's details.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the names of all caught pokemon.",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(args []string) error {
	commands := createCommandMap()

	if len(args) == 0 {
		fmt.Println("\nWelcome to the Pokedex!\n\nCommands:")
		for input, command := range commands {
			fmt.Printf("%s: %s\n", input, command.description)
		}
		return nil
	}

	command := commands[args[0]]
	fmt.Printf("%s: %s\n", command.name, command.description)
	return nil
}

func commandExit(args []string) error {
	fmt.Println("Go catch em all!")
	return nil
}

// MAP -------------------------------------------------------------
var mapMove = pokeapi.GetMap()

func commandMap(args []string) error {
	locations, err := mapMove(true)
	if err != nil {
		return errors.New("unable to get next poke-map")
	}
	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapBack(args []string) error {
	locations, err := mapMove(false)
	if err != nil {
		return errors.New("unable to get previous poke-map")
	}
	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

// EXPLORE -------------------------------------------------------------
var exploreLocation = pokeapi.Explore()

func commandExplore(args []string) error {
	if args == nil {
		return commandHelp([]string{"explore"})
	}
	if len(args) > 1 {
		return errors.New("explore accepts only one argument")
	}
	fmt.Println("Exploring", args[0], "...")
	pokemons, err := exploreLocation(args[0])
	if err != nil {
		return errors.New("unable to explore " + args[0])
	}
	fmt.Println("Found Pokemon:")
	for _, p := range pokemons {
		fmt.Println(" -", p.Name)
	}
	return nil
}

// CATCH -------------------------------------------------------------
var getPokedex = pokeapi.CreatePokedex()

func commandCatch(args []string) error {
	if args == nil {
		return commandHelp([]string{"catch"})
	}
	if len(args) > 1 {
		return errors.New("you can only catch one pokemon at a time")
	}

	success, err := getPokedex().Catch(args[0])
	if err != nil {
		return fmt.Errorf("unable to catch "+args[0], err)
	}
	fmt.Println("Threw a pokeball at", args[0])
	if success {
		fmt.Println(args[0], "was caught!")
	} else {
		fmt.Println(args[0], "escaped...")
	}
	return nil
}

// INSPECT -------------------------------------------------------------
func commandInspect(args []string) error {
	if args == nil {
		return commandHelp([]string{"inspect"})
	}

	for _, pkmn := range args {
		if details, ok := getPokedex()[pkmn]; !ok || !details.Caught {
			fmt.Println("You have not caught a", pkmn)
		} else {
			fmt.Println("Name:", details.Name)
			fmt.Println("ID:", details.Id)
			fmt.Println("Height:", details.Height)
			fmt.Println("Weight:", details.Weight)
			fmt.Println("Types:")
			for _, pkmnType := range details.Types {
				fmt.Println("	-", pkmnType.Type.Name)
			}
			fmt.Println("Stats:")
			for _, stats := range details.Stats {
				fmt.Printf("	- %s: %d\n", stats.Stat.Name, stats.BaseStat)
			}
			fmt.Println("Catch Rate:", details.CatchRate)
		}
	}
	return nil
}

// POKEDEX -------------------------------------------------------------
func commandPokedex(args []string) error {
	if len(getPokedex()) == 0 {
		fmt.Println("You have not caught any pokemon...")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for pkmn, details := range getPokedex() {
		status := "seen ✕"
		if details.Caught {
			status = "caught ✓"
		}
		fmt.Printf("	-%s (%s)\n", pkmn, status)
	}
	return nil
}
