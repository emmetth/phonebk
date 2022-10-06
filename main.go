package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emmetth/phonebk/contacts"
	_ "github.com/mattn/go-sqlite3"
)

var baseStyle = lipgloss.NewStyle()

func NewModel() (*model, error) {
	return &model{}, nil
}

type model struct {
	contacts []contacts.Contact
	cursor   int
	page     int
	pageSize int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) Home() {
	m.cursor = 0
}

func (m *model) End() {
	m.cursor = len(m.contacts) - 1
}

func (m *model) Up() {
	m.cursor = m.cursor - 1
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *model) Down() {
	m.cursor = m.cursor + 1
	if m.cursor > len(m.contacts)-1 {
		m.cursor = len(m.contacts) - 1
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "home", "alt+[H":
			m.Home()
		case "end", "alt+[F":
			m.End()
		case "down":
			m.Down()
		case "up":
			m.Up()
		}
	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	for i, c := range m.contacts {
		sb.WriteString(baseStyle.Reverse(i == m.cursor).Render(fmt.Sprintf("%-15s | %-15s | %-15s | %s", c.Lname, c.Fname, c.Phone, c.Email)) + "\n")
	}

	// details
	d := m.contacts[m.cursor]
	sb.WriteString(strings.Repeat("=", 80))
	sb.WriteString(fmt.Sprintf("\n%s %s, %s %s\n\n%s\n", d.Address, d.City, d.State, d.Zipcode, d.Notes))

	return sb.String()
}

func main() {
	conn, err := sql.Open("sqlite3", "phonebk.db")
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db := contacts.New(conn)

	contacts, err := db.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	m := model{contacts, 0, 0, 0}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
