package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MowlCoder/go-repl-pokedex/pokecache"
	"io"
	"net/http"
	"time"
)

var cache = pokecache.NewCache(pokecache.Config{
	Interval: 5 * time.Minute,
})

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(url string) (*LocationResponse, error) {
	var body []byte

	if bodyFromCache, ok := cache.Get(url); ok {
		body = bodyFromCache
	} else {
		res, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if res.StatusCode > 299 {
			return nil, errors.New("can not get locations from pokeapi")
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}

		cache.Add(url, body)
	}

	locationResponse := LocationResponse{}

	if err := json.Unmarshal(body, &locationResponse); err != nil {
		return nil, err
	}

	return &locationResponse, nil
}

type LocationAreaResponse struct {
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
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

func GetLocationArea(areaName string) (*LocationAreaResponse, error) {
	var body []byte
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	if bodyFromCache, ok := cache.Get(url); ok {
		body = bodyFromCache
	} else {
		res, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if res.StatusCode > 299 {
			return nil, errors.New("can not get location area from pokeapi")
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}

		cache.Add(url, body)
	}

	locationAreaResponse := LocationAreaResponse{}

	if err := json.Unmarshal(body, &locationAreaResponse); err != nil {
		return nil, err
	}

	return &locationAreaResponse, nil
}

type PokemonInfoResponse struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func GetPokemonInfo(pokemonName string) (*PokemonInfoResponse, error) {
	var body []byte
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	if bodyFromCache, ok := cache.Get(url); ok {
		body = bodyFromCache
	} else {
		res, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if res.StatusCode > 299 {
			return nil, errors.New("can not get pokemon info from pokeapi")
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}

		cache.Add(url, body)
	}

	pokemonInfoResponse := PokemonInfoResponse{}

	if err := json.Unmarshal(body, &pokemonInfoResponse); err != nil {
		return nil, err
	}

	return &pokemonInfoResponse, nil
}
