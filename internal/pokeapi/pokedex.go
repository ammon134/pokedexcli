package pokeapi

import "sync"

func PokedexInit() Pokedex {
	return Pokedex{
		pokedex: make(map[string]Pokemon),
		mu:      &sync.Mutex{},
	}
}

func (pd Pokedex) Add(pokemonData Pokemon) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	pd.pokedex[pokemonData.Name] = pokemonData
}
