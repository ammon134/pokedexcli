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
