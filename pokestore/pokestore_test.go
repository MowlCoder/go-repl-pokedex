package pokestore

import (
	"fmt"
	"testing"
)

func TestAddGet(t *testing.T) {
	cases := []struct {
		key     string
		pokemon Pokemon
	}{
		{
			key: "test-pokemon",
			pokemon: Pokemon{
				ID:   1,
				Name: "test-pokemon",
			},
		},
		{
			key: "test-pokemon2",
			pokemon: Pokemon{
				ID:   2,
				Name: "test-pokemon2",
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			store := NewPokemonStore()
			store.Add(c.pokemon)
			val, ok := store.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if val.ID != c.pokemon.ID {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}
