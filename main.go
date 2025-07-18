package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	callback func(p *Config, c *internal.Cache, flag string) error
}

type LocationArea struct {
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
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
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
func commandExit(p *Config, c *internal.Cache, flag string) error {
	fmt.Println("")
	
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)
	return nil
}

//called on "help" command
func commandHelp(p *Config, c *internal.Cache, flag string) error {
	fmt.Println("")

	commands := make(map[string]string)
	commands["exit"] = "Exit the Pokedex"
	commands["help"] = "Displays a help message"
	commands["map"] = "Retrieves a group of locations"
	commands["mapb"] = "Retrieves the previous group of locations"
	commands["explore"] = "Provides a list of pokemon in a given area"

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
func commandMap(p *Config, c *internal.Cache, flag string) error { //currently gets fucked if it reaches the end of the data presumably
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
		
		fmt.Println("")

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

	fmt.Println("")
	
	return nil
}

//called on "mapb" command
func commandMapB(p *Config, c *internal.Cache, flag string) error {
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
	
		fmt.Println("")
	
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

	fmt.Println("")

	return nil
}

func commandExplore(p *Config, c *internal.Cache, flag string) error {
	// gets data from location-area/{flag}, makes LocationArea struct, returns list of struct.Pokemon.Name 
	
	if flag == "" {
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
			lost := fmt.Sprintf("error - %s not a valid area to explore", flag)
			fmt.Println(lost)
			return nil
		}
	
		var area LocationArea
		var names []byte
		var namesSlice []string

		err = json.Unmarshal(bod, &area)
		if err != nil {
			fmt.Println("error -", err)
			return err
		}
		for _, v := range area.PokemonEncounters {
			fmt.Println(v.Pokemon.Name)
			namesSlice = append(namesSlice, v.Pokemon.Name)
		}
		joinedString := strings.Join(namesSlice, " ")
		names = []byte(joinedString)
		
		//add data to cache
		c.Add(flag, names)
	
		fmt.Println("")

	
		return nil
	}

	allNames := string(pokemon)
	nameSlice := strings.Split(allNames, " ")
	for _, name := range nameSlice {
		fmt.Println(name)
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
	}
	for ;; {
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
				comm.callback(&param, work, commFlag)
				found = true
			}
		}
		if !found {
			fmt.Println("Unknown command")
		}
		commFlag = ""
		//resp := fmt.Sprintf("Your command was: %s", clean[0])
		//fmt.Println(resp)
	}
}
