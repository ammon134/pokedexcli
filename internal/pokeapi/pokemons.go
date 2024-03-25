package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/ammon134/pokedexcli/internal/pokecache"
)

func CatchPokemon(pokemon *Pokemon, pokedex Pokedex) (bool, error) {
	catchProbability := (1 - float64(pokemon.BaseExperience)/100000.0) * 0.6
	// r := rand.New(rand.NewSource(1000))
	res := rand.Float64()

	if res <= catchProbability {
		pokedex.Add(*pokemon)
		return true, nil
	}
	return false, nil
}

func GetPokemonData(pokemonName string, cache pokecache.Cache) (*Pokemon, error) {
	pokemon := &Pokemon{}

	url := baseURL + "/pokemon/" + pokemonName

	if cacheVal, ok := cache.Get(url); ok {
		err := json.Unmarshal(cacheVal, pokemon)
		if err != nil {
			return nil, err
		}

		return pokemon, nil
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

	err = json.Unmarshal(body, pokemon)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}
