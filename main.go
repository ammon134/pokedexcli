package main

import (
	"time"

	"github.com/ammon134/pokedexcli/internal/pokecache"
)

func main() {
	ca := &cmdArg{
		nextLocations: nil,
		prevLocations: nil,
		cache:         pokecache.NewCache(5 * time.Minute),
	}
	startRepl(ca)
}
