package pokestore

import (
	"fmt"
	"strings"
	"sync"
)

type Pokemon struct {
	ID     int
	Name   string
	Height int
	Weight int
	Stats  map[string]int
	Types  []string
}

func (p Pokemon) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Name: %s\n", p.Name))
	sb.WriteString(fmt.Sprintf("Height: %d\n", p.Height))
	sb.WriteString(fmt.Sprintf("Weight: %d\n", p.Weight))
	sb.WriteString("Stats:\n")

	for statName, statVal := range p.Stats {
		sb.WriteString(fmt.Sprintf("  - %s: %d\n", statName, statVal))
	}

	sb.WriteString("Types:\n")

	for _, pokeType := range p.Types {
		sb.WriteString(fmt.Sprintf("  - %s\n", pokeType))
	}

	return strings.TrimRight(sb.String(), "\n")
}

type PokemonStore struct {
	store map[string]Pokemon
	mu    *sync.RWMutex
}

func NewPokemonStore() *PokemonStore {
	return &PokemonStore{
		store: map[string]Pokemon{},
		mu:    &sync.RWMutex{},
	}
}

func (p *PokemonStore) Add(pokemon Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.store[pokemon.Name] = pokemon
}

func (p *PokemonStore) Get(key string) (Pokemon, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	val, ok := p.store[key]

	return val, ok
}

func (p *PokemonStore) GetAllKeys() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	keys := make([]string, 0)

	for key, _ := range p.store {
		keys = append(keys, key)
	}

	return keys
}
