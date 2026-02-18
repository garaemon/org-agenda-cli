package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	inputView
)

type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeSideBySide
	ViewModeBoard
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
	sortBy      string
	sortDesc    bool
	textInput   textinput.Model
	inputMode   string // "schedule" or "deadline"
	viewMode    ViewMode
	width       int
	height      int
	boardLists  [3]list.Model
	focusedInt  int // 0, 1, 2 for board columns
	keys        KeyMap
}

func NewModel(items []*item.Item, start time.Time, viewRange string, title string, sortBy string, sortDesc bool) Model {

	ti := textinput.New()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 30

	m := Model{
		allItems:    items,
		currentDate: start,
		viewRange:   viewRange,
		title:       title,
		state:       listView,
		sortBy:      sortBy,
		sortDesc:    sortDesc,
		textInput:   ti,
		viewMode:    ViewModeList,
		keys:        NewKeyMap(),
	}
	// Initialize list with empty items, will be populated by refreshList
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return m.keys.ShortHelp()
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return m.keys.FullHelp()
	}
	m.list = l

	// Initialize board lists
	for i := 0; i < 3; i++ {
		m.boardLists[i] = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
		m.boardLists[i].SetShowHelp(false)
		m.boardLists[i].SetShowTitle(true)
		switch i {
		case 0:
			m.boardLists[i].Title = "TODO"
		case 1:
			m.boardLists[i].Title = "WAITING"
		case 2:
			m.boardLists[i].Title = "DONE"
		}
	}

	m.refreshList()

	return m
}

func (m *Model) refreshList() {
	var filtered []*item.Item

	if m.viewRange == "" {
		filtered = m.allItems
	} else {
		start := m.currentDate
		var end time.Time
		switch m.viewRange {
		case "week":
			end = start.AddDate(0, 0, 6)
		case "month":
			y, mm, _ := start.Date()
			start = time.Date(y, mm, 1, 0, 0, 0, 0, start.Location())
			end = start.AddDate(0, 1, 0)
		default:
			end = start
		}
		filtered = agenda.FilterItemsByRange(m.allItems, start, end)
	}

	agenda.SortItems(filtered, m.sortBy, m.sortDesc)

	var listItems []list.Item
	for _, it := range filtered {
		listItems = append(listItems, ListItem{Item: it})
	}

	m.list.SetItems(listItems)

	// Populate board lists
	var todoItems []list.Item
	var waitingItems []list.Item
	var doneItems []list.Item

	for _, it := range filtered {
		li := ListItem{Item: it}
		switch it.Status {
		case item.StatusTodo:
			todoItems = append(todoItems, li)
		case item.StatusWaiting:
			waitingItems = append(waitingItems, li)
		case item.StatusDone:
			doneItems = append(doneItems, li)
		default:
			// No status -> treat as TODO? or separate?
			// Treat as TODO for now
			todoItems = append(todoItems, li)
		}
	}

	m.boardLists[0].SetItems(todoItems)
	m.boardLists[1].SetItems(waitingItems)
	m.boardLists[2].SetItems(doneItems)

	sortStatus := ""
	if m.sortBy != "" {
		order := "Asc"
		if m.sortDesc {
			order = "Desc"
		}
		sortStatus = fmt.Sprintf(" [Sort: %s %s]", m.sortBy, order)
	}

	if m.viewRange == "" {
		m.list.Title = fmt.Sprintf("Agenda: All Items%s", sortStatus)
	} else {
		start := m.currentDate
		end := start
		switch m.viewRange {
		case "week":
			end = start.AddDate(0, 0, 6)
		case "month":
			y, mm, _ := start.Date()
			start = time.Date(y, mm, 1, 0, 0, 0, 0, start.Location())
			end = start.AddDate(0, 1, 0)
		default:
			end = start
		}
		m.list.Title = fmt.Sprintf("Agenda: %s - %s%s", start.Format("2006-01-02"), end.Format("2006-01-02"), sortStatus)
	}
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

		if key.Matches(msg, m.keys.ToggleView) {
			// Cycle view modes
			switch m.viewMode {
			case ViewModeList:
				m.viewMode = ViewModeSideBySide
			case ViewModeSideBySide:
				m.viewMode = ViewModeBoard
			case ViewModeBoard:
				m.viewMode = ViewModeList
			}

			// Re-layout immediately
			m.layout(m.width, m.height)
			return m, nil
		}
		switch m.state {
		case detailView:
			switch msg.String() {
			case "q", "esc", "backspace":
				m.state = listView
				return m, nil
			}
		case listView:
			if m.viewMode == ViewModeBoard {
				break
			}
			switch msg.String() {
			case "enter":
				i, ok := m.list.SelectedItem().(ListItem)
				if ok {
					m.state = detailView
					content := fmt.Sprintf("# %s\n\n%s", i.Item.Title, i.Item.RawContent)
					m.viewport.SetContent(content)
					return m, nil
				}
			}

			// Use key.Matches for our custom keys
			switch {
			case key.Matches(msg, m.keys.NextPeriod):
				switch m.viewRange {
				case "week":
					m.currentDate = m.currentDate.AddDate(0, 0, 7)
				case "month":
					y, mm, _ := m.currentDate.Date()
					m.currentDate = time.Date(y, mm, 1, 0, 0, 0, 0, m.currentDate.Location()).AddDate(0, 1, 0)
				default:
					m.currentDate = m.currentDate.AddDate(0, 0, 1)
				}
				m.refreshList()
				return m, nil
			case key.Matches(msg, m.keys.PrevPeriod):
				switch m.viewRange {
				case "week":
					m.currentDate = m.currentDate.AddDate(0, 0, -7)
				case "month":
					y, mm, _ := m.currentDate.Date()
					m.currentDate = time.Date(y, mm, 1, 0, 0, 0, 0, m.currentDate.Location()).AddDate(0, -1, 0)
				default:
					m.currentDate = m.currentDate.AddDate(0, 0, -1)
				}
				m.refreshList()
				return m, nil
			case key.Matches(msg, m.keys.Sort):
				// Cycle sort criteria: file -> priority -> date -> status -> file
				switch m.sortBy {
				case "":
					m.sortBy = "priority"
				case "priority":
					m.sortBy = "date"
				case "date":
					m.sortBy = "status"
				case "status":
					m.sortBy = ""
				}
				m.refreshList()
				return m, nil
			case key.Matches(msg, m.keys.SortOrder):
				m.sortDesc = !m.sortDesc
				m.refreshList()
				return m, nil
			case key.Matches(msg, m.keys.ToggleStatus):
				if i, ok := m.list.SelectedItem().(ListItem); ok {
					_ = ToggleStatus(i.Item)
					m.refreshList()
				}
				return m, nil
			case key.Matches(msg, m.keys.CyclePriority):
				if i, ok := m.list.SelectedItem().(ListItem); ok {
					_ = CyclePriority(i.Item)
					m.refreshList()
				}
				return m, nil
			case key.Matches(msg, m.keys.Schedule):
				m.inputMode = "SCHEDULED"
				m.textInput.Placeholder = "YYYY-MM-DD"
				m.textInput.SetValue("")
				m.textInput.Focus()
				m.state = inputView
				return m, nil
			case key.Matches(msg, m.keys.Deadline):
				m.inputMode = "DEADLINE"
				m.textInput.Placeholder = "YYYY-MM-DD"
				m.textInput.SetValue("")
				m.textInput.Focus()
				m.state = inputView
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout(m.width, m.height)
	}

	// Handle Board View updates specifically if active
	if m.viewMode == ViewModeBoard && m.state == listView {
		// Intercept navigation keys for columns
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.BoardLeft):
				m.focusedInt--
				if m.focusedInt < 0 {
					m.focusedInt = 0
				}
				return m, nil
			case key.Matches(msg, m.keys.BoardRight):
				m.focusedInt++
				if m.focusedInt > 2 {
					m.focusedInt = 2
				}
				return m, nil
			case key.Matches(msg, m.keys.ToggleStatus):
				// Toggle status of selected item in focused list
				if i, ok := m.boardLists[m.focusedInt].SelectedItem().(ListItem); ok {
					_ = ToggleStatus(i.Item)
					m.refreshList()
				}
				return m, nil
			}
		}

		// Delegate to focused list
		var cmd tea.Cmd
		m.boardLists[m.focusedInt], cmd = m.boardLists[m.focusedInt].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch m.state {
	case listView:
		var cmd tea.Cmd
		if m.viewMode != ViewModeBoard {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

		// Sync viewport content if side-by-side
		if m.viewMode == ViewModeSideBySide {
			if i, ok := m.list.SelectedItem().(ListItem); ok {
				content := fmt.Sprintf("# %s\n\n%s", i.Item.Title, i.Item.RawContent)
				m.viewport.SetContent(content)
			}
		}
	case detailView:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	case inputView:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				val := m.textInput.Value()
				if i, ok := m.list.SelectedItem().(ListItem); ok {
					if err := UpdateTimestamp(i.Item, m.inputMode, val); err != nil {
						// Stay in input view and maybe show error?
						// For now just logging to console or ignoring is bad UX but acceptable for Phase 1 MVP
						// Ideally textInput could show error style.
					} else {
						m.state = listView
						m.refreshList()
					}
				}
			case tea.KeyEsc:
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.state == inputView {
		return fmt.Sprintf(
			"Enter %s:\n\n%s\n\n(esc to cancel, enter to save)",
			m.inputMode,
			m.textInput.View(),
		)
	}
	if m.state == detailView {
		return fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View())
	}
	if m.viewMode == ViewModeSideBySide {
		// Update viewport content just in case? Or rely on Update?
		// Rely on Update.
		// Layout horizontal.
		// We use ViewSideBySide (assuming implementation matches)
		// But I put logic in side_view.go but commented out rendering logic there.
		// Wait, side_view.go does:
		// return lipgloss.JoinHorizontal(...)
		// It re-renders list view.
		return ViewSideBySide(m)
	}

	if m.viewMode == ViewModeBoard {
		return ViewBoard(m)
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

func (m *Model) layout(width, height int) {
	h, v := docStyle.GetFrameSize()

	switch m.viewMode {
	case ViewModeList:
		m.list.SetSize(width-h, height-v)
		headerHeight := lipgloss.Height(m.headerView())
		m.viewport = viewport.New(width, height-headerHeight)
		m.viewport.YPosition = headerHeight

	case ViewModeSideBySide:
		halfWidth := (width - h) / 2
		m.list.SetSize(halfWidth, height-v)

		headerHeight := lipgloss.Height(m.headerView())
		m.viewport = viewport.New(width-halfWidth, height-headerHeight)
		m.viewport.YPosition = headerHeight

		if i, ok := m.list.SelectedItem().(ListItem); ok {
			content := fmt.Sprintf("# %s\n\n%s", i.Item.Title, i.Item.RawContent)
			m.viewport.SetContent(content)
		}

	case ViewModeBoard:
		// 3 columns
		colWidth := (width - h) / 3
		// Adjust for padding/borders of columns?
		// board.go defines styles with padding/border.
		// We should account for that. Let's precise layout later.
		for i := 0; i < 3; i++ {
			m.boardLists[i].SetSize(colWidth-4, height-v-2) // rough estimate
		}
	}
}
