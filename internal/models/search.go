package models

import (
	"autocomplete/internal/models/trie"

	"github.com/charmbracelet/bubbles/textinput"
)

type SearchModel struct {
	TextInput   textinput.Model
	Service     trie.AutocompleteService
	Suggestions []trie.SuggestionResult
	Selected    int
	Error       error
}

func NewSearchModel(config trie.Config) (*SearchModel, error) {
	service, err := trie.NewAutocompleteService(config)
	if err != nil {
		return nil, err
	}

	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 50

	return &SearchModel{
		TextInput:   ti,
		Service:     service,
		Suggestions: []trie.SuggestionResult{},
		Selected:    0,
	}, nil
}

func (m *SearchModel) UpdateSuggestions(input string, limit int) error {
	if len(input) == 0 {
		m.Suggestions = []trie.SuggestionResult{}
		return nil
	}

	suggestions, err := m.Service.Search(input, limit)
	if err != nil {
		return err
	}

	m.Suggestions = suggestions
	m.Selected = 0
	return nil
}
