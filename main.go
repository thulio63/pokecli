package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

//interface for any commands we make
type cliCommand struct {
	name string
	description string
	callback func(p *Config) error
}

type Result struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type Config struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []Result `json:"results"` //{
		//Name string `json:"name"`
		//URL string `json:"url"`
	//} `json:"results"`
}

//ADD NEW COMMANDS TO MAP IN MAIN()

//called on "exit" command
func commandExit(p *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

//called on "help" command
func commandHelp(p *Config) error {
	commands := make(map[string]string)
	commands["exit"] = "Exit the Pokedex"
	commands["help"] = "Displays a help message"
	commands["map"] = "Retrieves a group of locations"
	commands["mapb"] = "Retrieves the previous group of locations"
	//commandNames := []string{"exit", "help"}
	//commandDescriptions := []string{"Exit the Pokedex", "Displays a help message"}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for comm, scrip := range commands {
		resp := fmt.Sprintf("%s: %s", comm, scrip)
		fmt.Println(resp)
	}
	fmt.Println("")
	return nil
}

//called on "map" command
func commandMap(p *Config) error {
	resp, err := http.Get(p.Next)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	con := Config{}
	err = json.Unmarshal(bod, &con)
	if err != nil {
		return err
	}
	p.Next = con.Next
	p.Previous = con.Previous
	
	for _, res := range con.Results {
		fmt.Println(res.Name)
	}

	fmt.Println("")

	return nil
}

//called on "mapb" command
func commandMapB(p *Config) error {
	if p.Previous == "null" {
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := http.Get(p.Previous)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	con := Config{}
	err = json.Unmarshal(bod, &con)
	if err != nil {
		return err
	}
	p.Next = con.Next
	p.Previous = con.Previous
	
	for _, res := range con.Results {
		fmt.Println(res.Name)
	}

	fmt.Println("")

	return nil
}

func cleanInput(text string) []string {
	mySlice := strings.Fields(strings.ToLower(text))
	return mySlice
}



func main() {
	input := bufio.NewScanner(os.Stdin)
	param := Config{0, "https://pokeapi.co/api/v2/location-area/", "null", []Result {},}
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
		"map": {
			name: "map",
			description: "Retrieves a group of locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Retrieves the previous group of locations",
			callback: commandMapB,
		},
	}
	for ;; {
		fmt.Print("Pokedex > ")
		input.Scan()
		command := cleanInput(input.Text())[0]
		found := false
		for _, comm := range commandList {
			if command == comm.name {
				comm.callback(&param)
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
