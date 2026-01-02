package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

type ListItem struct {
	Item *item.Item
}

func (i ListItem) Title() string {
	return i.Item.Title
}

func (i ListItem) Description() string {
	var parts []string
	if i.Item.Status != "" {
		parts = append(parts, fmt.Sprintf("[%s]", i.Item.Status))
	}
	if i.Item.Scheduled != nil {
		parts = append(parts, fmt.Sprintf("Sch: %s", i.Item.Scheduled.Format("2006-01-02")))
	}
	if i.Item.Deadline != nil {
		parts = append(parts, fmt.Sprintf("Ddl: %s", i.Item.Deadline.Format("2006-01-02")))
	}
	if len(i.Item.Tags) > 0 {
		parts = append(parts, fmt.Sprintf(":%s:", strings.Join(i.Item.Tags, ":")))
	}
	return strings.Join(parts, " ")
}

func (i ListItem) FilterValue() string {
	return i.Item.Title
}

type Model struct {
	list list.Model
}

func NewModel(items []*item.Item, title string) Model {
	var listItems []list.Item
	for _, it := range items {
		listItems = append(listItems, ListItem{Item: it})
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = title

	return Model{list: l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}
