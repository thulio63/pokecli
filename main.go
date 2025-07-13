package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//interface for any commands we make
type cliCommand struct {
	name string
	description string
	callback func() error
}

//ADD NEW COMMANDS TO MAP IN MAIN()

//called on "exit" command
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

//called on "help" command
func commandHelp() error {
	commandNames := []string{"exit", "help"}
	commandDescriptions := []string{"Exit the Pokedex", "Displays a help message"}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for i, comm := range commandNames {
		resp := fmt.Sprintf("%s: %s", comm, commandDescriptions[i])
		fmt.Println(resp)
	}
	return nil
}

func cleanInput(text string) []string {
	mySlice := strings.Fields(strings.ToLower(text))
	return mySlice
}



func main() {
	input := bufio.NewScanner(os.Stdin)

	commandList := map[string]cliCommand{
	
		"exit": {
			name: "exit",
			description: "Exit the Pokedex", 
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
	}
	for ;; {
		fmt.Print("Pokedex > ")
		input.Scan()
		command := cleanInput(input.Text())[0]
		found := false
		for _, comm := range commandList {
			if command == comm.name {
				comm.callback()
				found = true
			}
		}
		if !found {
			fmt.Println("Unknown command")
		}
		//resp := fmt.Sprintf("Your command was: %s", clean[0])
		//fmt.Println(resp)
	}
}
