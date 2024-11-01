package controllers

import (
	"autocomplete/internal/models"
	"autocomplete/internal/models/trie"
	"autocomplete/internal/views"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchController struct {
	model *models.SearchModel
	view  *views.SearchView
}

func NewSearchController(model *models.SearchModel) *SearchController {
	return &SearchController{
		model: model,
		view:  views.NewSearchView(model),
	}
}

func (c *SearchController) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		c.updateSuggestions,
	)
}

func (c *SearchController) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, tea.Quit
		case tea.KeyEnter:
			return c.handleEnter()
		case tea.KeyUp:
			return c.handleKeyUp()
		case tea.KeyDown:
			return c.handleKeyDown()
		}
	}

	c.model.TextInput, cmd = c.model.TextInput.Update(msg)

	currentInput := c.model.TextInput.Value()
	var err error
	if len(currentInput) > 0 {
		c.model.Suggestions, err = c.model.Service.Search(currentInput, 5)
		if err != nil {
			fmt.Println("Error in searching", err)
		}
		c.model.Selected = 0 // Reset selection when suggestions change
	} else {
		c.model.Suggestions = []trie.SuggestionResult{}
	}
	return c, tea.Batch(cmd, c.updateSuggestions)
}

func (c *SearchController) View() string {
	return c.view.Render()
}

func (c *SearchController) handleEnter() (tea.Model, tea.Cmd) {
	if len(c.model.Suggestions) > 0 {
		selected := c.model.Suggestions[c.model.Selected]
		if err := c.model.Service.RecordSelection(selected.Word); err != nil {
			c.model.Error = err
			return c, nil
		}
		c.model.TextInput.SetValue(selected.Word)
		c.model.Suggestions = []trie.SuggestionResult{}
	}
	return c, nil
}

func (c *SearchController) handleKeyUp() (tea.Model, tea.Cmd) {
	if c.model.Selected > 0 {
		c.model.Selected--
	}
	return c, nil
}

func (c *SearchController) handleKeyDown() (tea.Model, tea.Cmd) {
	if c.model.Selected < len(c.model.Suggestions)-1 {
		c.model.Selected++
	}
	return c, nil
}

func (c *SearchController) updateSuggestions() tea.Msg {
	currentInput := c.model.TextInput.Value()
	c.model.UpdateSuggestions(currentInput, 5)
	return nil
}
