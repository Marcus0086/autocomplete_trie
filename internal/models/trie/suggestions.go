package trie

import "time"

type SuggestionResult struct {
	Word      string    `json:"word,omitempty"`
	Score     float64   `json:"score,omitempty"`
	Frequency int       `json:"frequency,omitempty"`
	LastUsed  time.Time `json:"lastUsed,omitempty"`
}
