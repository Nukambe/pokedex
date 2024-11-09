package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

type locationResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationArea struct {
	Results  []locationResult `json:"results"`
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
}

type pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type exploreResult struct {
	Pokemon pokemon `json:"pokemon"`
}

type exploreArea struct {
	PokemonEncounters []exploreResult `json:"pokemon_encounters"`
}

type stat struct {
	Stat     pokemon `json:"stat"`
	BaseStat int     `json:"base_stat"`
	Effort   int     `json:"effort"`
}

type pkmnType struct {
	Type pokemon `json:"type"`
	Slot int     `json:"slot"`
}

type species struct {
	BaseHappiness int `json:"base_happiness"`
	CaptureRate   int `json:"capture_rate"`
}

type pokemonDetails struct {
	Stats          []stat     `json:"stats"`
	Types          []pkmnType `json:"types"`
	Species        pokemon    `json:"species"`
	Name           string     `json:"name"`
	Id             int        `json:"id"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	BaseExperience int        `json:"base_experience"`
	CatchRate      int
}

type pokedex map[string]pokemonDetails

func GetMap() func(bool) ([]locationResult, error) {
	currentUrl := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	nextUrl := ""
	previousUrl := ""
	url := &currentUrl
	locations := map[string]locationArea{}

	return func(next bool) ([]locationResult, error) {
		if next && nextUrl != "" {
			url = &nextUrl
		} else if !next && previousUrl != "" {
			url = &previousUrl
		} else {
			url = &currentUrl
		}

		if _, ok := locations[*url]; !ok {
			res, err := http.Get(*url)
			if err != nil {
				return nil, errors.New("network error fetching pokeapi locations")
			}
			defer res.Body.Close()

			var location locationArea
			decoder := json.NewDecoder(res.Body)
			if err = decoder.Decode(&location); err != nil {
				return nil, errors.New("json error decoding pokeapi locations")
			}
			locations[*url] = location
		}

		currentUrl = *url
		if locations[currentUrl].Next == nil {
			nextUrl = ""
		} else {
			nextUrl = *locations[currentUrl].Next
		}
		if locations[currentUrl].Previous == nil {
			previousUrl = ""
		} else {
			previousUrl = *locations[currentUrl].Previous
		}

		return locations[currentUrl].Results, nil
	}
}

func Explore() func(string) ([]pokemon, error) {
	locations := map[string]exploreArea{}

	return func(location string) ([]pokemon, error) {
		if _, ok := locations[location]; !ok {
			res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + location)
			if err != nil {
				return nil, errors.New("network error exploring pokeapi location " + location)
			}
			defer res.Body.Close()

			var area exploreArea
			decoder := json.NewDecoder(res.Body)
			if err = decoder.Decode(&area); err != nil {
				return nil, errors.New("json error decoding pokeapi exploration " + location)
			}

			locations[location] = area
		}
		pokemons := make([]pokemon, len(locations[location].PokemonEncounters))
		for i, p := range locations[location].PokemonEncounters {
			pokemons[i] = p.Pokemon
		}
		return pokemons, nil
	}
}

func CreatePokedex() func() pokedex {
	pkdx := pokedex{}
	return func() pokedex {
		return pkdx
	}
}

func (pkdx pokedex) Catch(pkmn string) (success bool, err error) {
	if _, ok := pkdx[pkmn]; !ok {
		// GET POKEMON DETAILS
		res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pkmn)
		if err != nil {
			return false, errors.New("network error getting pokemon")
		}
		defer res.Body.Close()
		var details pokemonDetails
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&details); err != nil {
			return false, errors.New("json error decoding pokemon details")
		}
		// FROM POKEMON DETAILS, GET SPECIES FOR CATCH RATE
		resp, err := http.Get(details.Species.Url)
		if err != nil {
			return false, errors.New("network error getting species")
		}
		defer resp.Body.Close()
		var spcs species
		decoder = json.NewDecoder(resp.Body)
		if err := decoder.Decode(&spcs); err != nil {
			return false, fmt.Errorf("json error decoding species: %w", err)
		}
		// ADD CATCH RATE TO POKEMON DETAILS
		details.CatchRate = spcs.CaptureRate
		// SAVE POKEMON DETAILS TO THE POKEDEX
		pkdx[pkmn] = details
	}

	// CATCH RATE IS 0-255
	n := rand.Intn(256)
	return n <= pkdx[pkmn].CatchRate, nil
}
