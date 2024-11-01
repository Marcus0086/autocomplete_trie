package trie

import (
	"fmt"
	"time"
)

type trieService struct {
	trie   *Trie
	config Config
}

func NewAutocompleteService(config Config) (AutocompleteService, error) {
	trie, err := NewTrie(config.CacheSize, config.MaxResults)
	if err != nil {
		return nil, fmt.Errorf("error creating trie: %w", err)
	}

	words := map[float64]string{
		1.0: "elastic",
		0.9: "elasticsearch",
		0.8: "elasticity",
		0.7: "elbow",
		0.6: "element",
		0.5: "elementary",
		0.4: "elevator",
		0.3: "elephant",
		0.2: "elevation",
		0.1: "elegant",
	}

	for score, word := range words {
		trie.Insert(word, score)
	}
	return &trieService{
		trie:   trie,
		config: config,
	}, nil
}

func (s *trieService) Search(prefix string, limit int) ([]SuggestionResult, error) {
	results := s.trie.SearchWithSuggestions(prefix, limit)
	return results, nil
}

func (s *trieService) AddWord(word string, score float64) error {
	s.trie.Insert(word, score)
	return nil
}

func (s *trieService) RecordSelection(word string) error {
	s.trie.RecordSelection(word)
	return nil
}

func (s *trieService) Cleanup(maxAge time.Duration) error {
	s.trie.Cleanup(maxAge)
	return nil
}
