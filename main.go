package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
	clipboard2 "myclipboard/clipboard"
	"myclipboard/storage/file"
	"os"
	"time"
)

type Model struct {
	choices  list.Model
	storage  *file.Storage
	cursor   int
	selector int
}

type tickMsg struct{}

func ticker() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func initialModel(s *file.Storage) Model {
	_ = s.Load(20)
	clipboard2.InitClipboard(s)
	return Model{
		choices:  list.New(s.ToBubbleList(), list.NewDefaultDelegate(), 0, 0),
		selector: -1,
		storage:  s,
	}
}

func (m Model) Init() tea.Cmd {
	return ticker()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices.Items())-1 {
				m.cursor++
			}
		case "enter", " ":
			m.storage.ClipboardSig <- struct{}{}
			m.selector = m.cursor
			m.choices.Select(m.cursor)
			clipboard.Write(clipboard.FmtText, []byte(m.choices.SelectedItem().FilterValue()))
			_ = m.storage.Select(m.cursor)
			m.choices.SetItems(m.storage.ToBubbleList())
			return m, nil
		}
	case tickMsg:
		m.choices.SetItems(m.storage.ToBubbleList())
		m.selector = -1
		return m, ticker()
	}
	return m, nil
}

func (m Model) View() string {
	s := ""
	for i, choice := range m.choices.Items() {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := "      "
		if m.cursor == m.selector && m.selector == i {
			checked = "copied"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	return s
}

func main() {
	s := file.Storage{}
	err := s.Init("./storage/data")
	if err != nil {
		return
	}
	defer s.Close()

	p := tea.NewProgram(initialModel(&s))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
