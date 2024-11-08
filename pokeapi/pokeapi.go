package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type locationResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationArea struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []locationResult `json:"results"`
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
