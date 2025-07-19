package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thulio63/pokecli/internal"
)

//interface for any commands we make
type cliCommand struct {
	name string
	description string
	callback func(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error
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

//ADD NEW COMMANDS TO MAP IN MAIN() AND IN COMMANDHELP

//called on "exit" command
func commandExit(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	fmt.Println("")

	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)
	return nil
}

//called on "help" command
func commandHelp(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	fmt.Println("")

	commands := make(map[string]string)
	commands["exit"] = "Exit the Pokedex"
	commands["help"] = "Displays a help message"
	commands["map"] = "Retrieves a group of locations"
	commands["mapb"] = "Retrieves the previous group of locations"
	commands["explore"] = "Provides a list of pokemon in a given area"
	commands["catch"] = "Attempts to catch a pokemon"
	commands["inspect"] = "View information on pokemon you have previously caught"
	commands["pokedex"] = "List pokemon you have previously caught"

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for comm, scrip := range commands {
		resp := fmt.Sprintf("%s: %s", comm, scrip)
		fmt.Println(resp)
	}

	return nil
}

//called on "map" command
func commandMap(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error { //currently gets fucked if it reaches the end of the data presumably
	fmt.Println("")

	//check if p.Next is in cache
	page, found := c.Get(p.Next)
	//if not, run http request and add data
	if !found {
		resp, err := http.Get(p.Next)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		defer resp.Body.Close()
		
		bod, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		con := Config{}
		err = json.Unmarshal(bod, &con)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		
		//add data to cache
		c.Add(p.Next, bod)

		p.Next = con.Next
		p.Previous = con.Previous
		
		for _, res := range con.Results {
			fmt.Println(res.Name)
		}

		return nil
	}
	//if yes, skip request and read from data
	bod := page
	con := Config{}
	err := json.Unmarshal(bod, &con)
	if err != nil {
		fmt.Println("error -", err)
		return err
	}
	p.Next = con.Next
	p.Previous = con.Previous
	
	for _, res := range con.Results {
		fmt.Println(res.Name)
	}
	
	return nil
}

//called on "mapb" command
func commandMapB(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	fmt.Println("")

	//returns if there is no previous page
	if p.Previous == "" {
		fmt.Println("you're on the first page, dumbass")
		return nil
	}

	//checks if previous page is in cache
	page, found := c.Get(p.Next)

	//if no, run http request and add
	if !found {
		resp, err := http.Get(p.Previous)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		defer resp.Body.Close()
	
		bod, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
	
		con := Config{}
		err = json.Unmarshal(bod, &con)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}

		//add data to cache
		c.Add(p.Previous, bod)

		p.Next = con.Next
		p.Previous = con.Previous
		
		for _, res := range con.Results {
			fmt.Println(res.Name)
		}
	
		return nil
	}
	//if yes, skip request and read from data
	bod := page

	con := Config{}
	err := json.Unmarshal(bod, &con)
	if err != nil {
		fmt.Println("error -", err)
		return err
	}
	p.Next = con.Next
	p.Previous = con.Previous
	
	for _, res := range con.Results {
		fmt.Println(res.Name)
	}

	return nil
}

func commandExplore(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	// gets data from location-area/{flag}, makes LocationArea struct, returns list of struct.Pokemon.Name 
	
	if flag == "" {
		fmt.Println("")
		fmt.Println("error - please provide an area to explore")
		return nil
	}

	//checks if location-area data is in cache
	pokemon, found := c.Get(flag)

	//if no, run http request and add
	if !found {
		url := "https://pokeapi.co/api/v2/location-area/" + flag 
		resp, err := http.Get(url)
		
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		defer resp.Body.Close()
	
		bod, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}

		if string(bod) == "Not Found" {
			fmt.Println("")
			lost := fmt.Sprintf("error - %s not a valid area to explore", flag)
			fmt.Println(lost)
			return nil
		}
	
		var area internal.LocationArea
		var names []byte
		var namesSlice []string

		err = json.Unmarshal(bod, &area)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		fmt.Println("")
		for _, v := range area.PokemonEncounters {
			fmt.Println(v.Pokemon.Name)
			namesSlice = append(namesSlice, v.Pokemon.Name)
		}
		joinedString := strings.Join(namesSlice, " ")
		names = []byte(joinedString)
		
		//add data to cache
		c.Add(flag, names)	

		return nil
	}

	allNames := string(pokemon)
	nameSlice := strings.Split(allNames, " ")
	for _, name := range nameSlice {
		fmt.Println(name)
	}

	return nil
}

func commandCatch(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	//checks for flag
	if flag == "" {
		fmt.Println("error - please attempt to catch a pokemon")
		return nil
	}
	
	fmt.Printf("\nThrowing a Pokeball at %s...\n", flag)
	
	//set odds for capture
	odds := 100
	success := 25

	//check if pokemon is in pokedex
	poke, ok := pd.Caught[flag]

	//if no, run http request for pokemon data
	if !ok {
		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + flag)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		defer resp.Body.Close()

		bod, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}

		if string(bod) == "Not Found" {
			fmt.Println("")
			fmt.Printf("error - %s not a valid pokemon\n", flag)
			return nil
		}

		var target internal.Pokemon

		err = json.Unmarshal(bod, &target)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}

		sub := target.Base_exp / 100

		if rand.IntN(odds) - sub > success {
			fmt.Println(flag, "was caught!")
			//make pokemon visible to user
			target.Captured = true
			fmt.Println(flag, "was added to the pokedex!")
			
		} else {
			fmt.Println(flag, "escaped!")
		}

		//add pokemon to pokedex with correct captured value 
		pd.Caught[flag] = target

		return nil
	} else { //pokemon in pokedex
		sub := poke.Base_exp / 100
		if rand.IntN(odds) - sub > success {
			fmt.Println(flag, "was caught!")
			if !poke.Captured {
				//set pokemon captured status to true
				update := pd.Caught[flag]
				update.Captured = true
				pd.Caught[flag] = update
				fmt.Println(flag, "was added to the pokedex!")
			}
			poke.Captured = true 
		} else {
			fmt.Println(flag, "escaped!")
		}
		return nil
	} 
}

func commandInspect(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	//checks for flag
	if flag == "" {
		fmt.Println("error - please attempt to inspect a pokemon")
		return nil
	}

	//check if pokemon is in pokedex
	poke, ok := pd.Caught[flag]

	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Println("Name:", poke.Name)
	fmt.Println("Height:", poke.Height)
	fmt.Println("Weight:", poke.Weight)
	fmt.Println("Stats:")

	for _, val := range poke.Stats {
		fmt.Printf(" -%s: %d\n", val.Stat.Name, val.BaseStat)
	}
	
	fmt.Println("Types:")

	for _, myType := range poke.Types {
		fmt.Println(" -", myType.Type.Name)
	}

	return nil
}

func commandPokedex(p *Config, c *internal.Cache, flag string, pd *internal.PokeDex) error {
	fmt.Println("Your Pokedex:")
	fmt.Println("")
	empty := true
	for _, caught := range pd.Caught {
		if caught.Captured {
			empty = false
			fmt.Println(" -", caught.Name)
		}
	}
	if empty {
		fmt.Println("You currently posess no Pokemon. Go catch 'em all!")
	}
	return nil
}

func cleanInput(text string) []string {
	mySlice := strings.Fields(strings.ToLower(text))
	return mySlice
}



func main() {
	empty := make(map[string]internal.Pokemon)
	pokedex := internal.PokeDex{User: "Andrew", Caught: empty}
	input := bufio.NewScanner(os.Stdin)
	param := Config{0, "https://pokeapi.co/api/v2/location-area/", "", []Result {},}
	cache := internal.NewCache(time.Second * 15)
	work := &cache
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
		"explore": {
			name: "explore",
			description: "Provides a list of pokemon in a given area",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "Attempts to catch a pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "View information on pokemon you have previously caught",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "List pokemon you have previously caught",
			callback: commandPokedex,
		},
	}
	for ;; {
		fmt.Println("")
		fmt.Print("Pokedex > ")
		input.Scan()
		commands := cleanInput(input.Text())
		//fmt.Println(len(commands))
		command := commands[0]
		commFlag := ""
		if len(commands) > 1 {
			commFlag = commands[1]
		} 
		found := false
		for _, comm := range commandList {
			if command == comm.name {
				comm.callback(&param, work, commFlag, &pokedex)
				found = true
			}
		}
		if !found {
			fmt.Println("")
			fmt.Println("Unknown command")
		}
		commFlag = ""
		//resp := fmt.Sprintf("Your command was: %s", clean[0])
		//fmt.Println(resp)
	}
}
