package trie

import (
	"time"
)

type AutocompleteService interface {
	Search(prefix string, limit int) ([]SuggestionResult, error)
	AddWord(word string, score float64) error
	RecordSelection(word string) error
	Cleanup(maxAge time.Duration) error
}

// Config holds the configuration for the autocomplete service
type Config struct {
	CacheSize    int
	MaxResults   int
	CleanupAge   time.Duration
	ScoreWeights ScoreWeights
}

// ScoreWeights defines weights for different scoring factors
type ScoreWeights struct {
	BaseScore     float64
	FrequencyMult float64
	RecencyMult   float64
}

// NewDefaultConfig returns default configuration
func NewDefaultConfig() Config {
	return Config{
		CacheSize:  1000,
		MaxResults: 10,
		CleanupAge: 24 * time.Hour,
		ScoreWeights: ScoreWeights{
			BaseScore:     1.0,
			FrequencyMult: 0.3,
			RecencyMult:   0.2,
		},
	}
}
