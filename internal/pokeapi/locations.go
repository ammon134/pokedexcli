package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ammon134/pokedexcli/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type LocationListResult struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type LocationDetailResult struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon        Pokemon `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocations(url *string, cache pokecache.Cache) (*LocationListResult, error) {
	locationRes := &LocationListResult{}
	getURL := baseURL + "/location-area?offset=0&limit=20"
	if url != nil {
		getURL = *url
	}

	if cacheVal, ok := cache.Get(getURL); ok {
		// fmt.Println("cache hit")
		err := json.Unmarshal(cacheVal, locationRes)
		if err != nil {
			return locationRes, err
		}
	}

	res, err := http.Get(getURL)
	if err != nil {
		return locationRes, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return locationRes, fmt.Errorf("failed with status code %d and\nbody: %s", res.StatusCode, body)
	}

	if err != nil {
		return locationRes, err
	}

	cache.Add(getURL, body)

	err = json.Unmarshal(body, locationRes)
	if err != nil {
		return locationRes, err
	}

	return locationRes, nil
}

func GetLocationDetail(location string, cache pokecache.Cache) (*LocationDetailResult, error) {
	lDRes := &LocationDetailResult{}

	url := baseURL + "/location-area/" + location

	if cacheVal, ok := cache.Get(url); ok {
		err := json.Unmarshal(cacheVal, lDRes)
		if err != nil {
			return nil, err
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("failed with status code %d and \nbody: %s", res.StatusCode, body)
	}

	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, lDRes)
	if err != nil {
		return nil, err
	}

	return lDRes, nil
}
