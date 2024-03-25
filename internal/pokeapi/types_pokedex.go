package pokeapi

import "sync"

type Pokedex struct {
	pokedex map[string]Pokemon
	mu      *sync.Mutex
}
