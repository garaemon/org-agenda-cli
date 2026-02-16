package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	NextPeriod    key.Binding
	PrevPeriod    key.Binding
	Sort          key.Binding
	SortOrder     key.Binding
	ToggleStatus  key.Binding
	CyclePriority key.Binding
	Schedule      key.Binding
	Deadline      key.Binding
	ToggleView    key.Binding
	BoardLeft     key.Binding
	BoardRight    key.Binding
	InputEnter    key.Binding
	InputEsc      key.Binding
	Quit          key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		NextPeriod: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next period"),
		),
		PrevPeriod: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "prev period"),
		),
		Sort: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "cycle sort"),
		),
		SortOrder: key.NewBinding(
			key.WithKeys("O"),
			key.WithHelp("O", "toggle sort order"),
		),
		ToggleStatus: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle status"),
		),
		CyclePriority: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "cycle priority"),
		),
		Schedule: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "set schedule"),
		),
		Deadline: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "set deadline"),
		),
		ToggleView: key.NewBinding(
			key.WithKeys("tab", "|"),
			key.WithHelp("tab", "toggle view"),
		),
		BoardLeft: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h/left", "board left"),
		),
		BoardRight: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l/right", "board right"),
		),
		InputEnter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		InputEsc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ToggleView, k.Sort, k.ToggleStatus, k.Schedule, k.Deadline, k.NextPeriod, k.PrevPeriod,
	}
}

func (k KeyMap) FullHelp() []key.Binding {
	return []key.Binding{
		k.NextPeriod, k.PrevPeriod,
		k.Sort, k.SortOrder,
		k.ToggleView,
		k.ToggleStatus, k.CyclePriority,
		k.Schedule, k.Deadline,
		k.BoardLeft, k.BoardRight,
		k.Quit,
	}
}
