package main

import (
	"autocomplete/internal/controllers"
	"autocomplete/internal/models"
	"autocomplete/internal/models/trie"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := trie.NewDefaultConfig()
	model, err := models.NewSearchModel(config)
	if err != nil {
		fmt.Printf("Error initializing model: %v\n", err)
		return
	}

	controller := controllers.NewSearchController(model)
	p := tea.NewProgram(controller)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return
	}
}
