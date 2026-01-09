package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garaemon/org-agenda-cli/pkg/agenda"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

type sessionState int

const (
	listView sessionState = iota
	detailView
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
	if i.Item.FilePath != "" {
		parts = append(parts, fmt.Sprintf("(%s:%d)", i.Item.FilePath, i.Item.LineNumber))
	}
	return strings.Join(parts, " ")
}

func (i ListItem) FilterValue() string {
	return i.Item.Title
}

type Model struct {
	list        list.Model
	viewport    viewport.Model
	state       sessionState
	allItems    []*item.Item
	currentDate time.Time
	viewRange   string
	title       string
}

func NewModel(items []*item.Item, start time.Time, viewRange string, title string) Model {
	m := Model{
		allItems:    items,
		currentDate: start,
		viewRange:   viewRange,
		title:       title,
		state:       listView,
	}
	// Initialize list with empty items, will be populated by refreshList
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list = l
	m.refreshList()

	return m
}

func (m *Model) refreshList() {
	if m.viewRange == "" {
		var listItems []list.Item
		for _, it := range m.allItems {
			listItems = append(listItems, ListItem{Item: it})
		}
		m.list.SetItems(listItems)
		m.list.Title = m.title
		return
	}

	start := m.currentDate
	end := start
	if m.viewRange == "week" {
		end = start.AddDate(0, 0, 7)
	} else if m.viewRange == "month" {
		end = start.AddDate(0, 1, 0)
	} else {
		// Default to day
		end = start
	}

	filtered := agenda.FilterItemsByRange(m.allItems, start, end)

	var listItems []list.Item
	for _, it := range filtered {
		listItems = append(listItems, ListItem{Item: it})
	}

	m.list.SetItems(listItems)
	m.list.Title = fmt.Sprintf("Agenda: %s - %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.state {
		case detailView:
			switch msg.String() {
			case "q", "esc", "backspace":
				m.state = listView
				return m, nil
			}
		case listView:
			switch msg.String() {
			case "enter":
				i, ok := m.list.SelectedItem().(ListItem)
				if ok {
					m.state = detailView
					content := fmt.Sprintf("# %s\n\n%s", i.Item.Title, i.Item.RawContent)
					m.viewport.SetContent(content)
					return m, nil
				}
			case "n":
				if m.viewRange == "week" {
					m.currentDate = m.currentDate.AddDate(0, 0, 7)
				} else if m.viewRange == "month" {
					m.currentDate = m.currentDate.AddDate(0, 1, 0)
				} else {
					m.currentDate = m.currentDate.AddDate(0, 0, 1)
				}
				m.refreshList()
				return m, nil
			case "p":
				if m.viewRange == "week" {
					m.currentDate = m.currentDate.AddDate(0, 0, -7)
				} else if m.viewRange == "month" {
					m.currentDate = m.currentDate.AddDate(0, -1, 0)
				} else {
					m.currentDate = m.currentDate.AddDate(0, 0, -1)
				}
				m.refreshList()
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		headerHeight := lipgloss.Height(m.headerView())
		m.viewport = viewport.New(msg.Width, msg.Height-headerHeight)
		m.viewport.YPosition = headerHeight
	}

	switch m.state {
	case listView:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case detailView:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.state == detailView {
		return fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View())
	}
	return docStyle.Render(m.list.View())
}

func (m Model) headerView() string {
	title := "Details"
	line := strings.Repeat("─", max(0, m.viewport.Width-len(title)))
	return lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(title + line)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
