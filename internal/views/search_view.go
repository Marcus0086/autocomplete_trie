package views

import (
	"autocomplete/internal/models"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	suggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Margin(0, 1)
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("36")).
			Bold(true)
)

type SearchView struct {
	model *models.SearchModel
}

func NewSearchView(model *models.SearchModel) *SearchView {
	return &SearchView{
		model: model,
	}
}

func (v *SearchView) Render() string {
	var s strings.Builder

	s.WriteString("Enter text to search:\n")
	s.WriteString(v.model.TextInput.View())
	s.WriteString("\n\n")

	if len(v.model.Suggestions) > 0 {
		s.WriteString("Suggestions:\n")
		for i, sugg := range v.model.Suggestions {
			if i == v.model.Selected {
				s.WriteString(selectedStyle.Render(fmt.Sprintf("▸ %s (%.2f)\n", sugg.Word, sugg.Score)))
			} else {
				s.WriteString(suggestionStyle.Render(fmt.Sprintf("  %s (%.2f)\n", sugg.Word, sugg.Score)))
			}
		}
	}

	return s.String()
}
