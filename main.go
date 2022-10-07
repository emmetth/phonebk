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
var db *contacts.Queries

const (
	StatusList = iota
	StatusConfirm
)

type model struct {
	contacts []contacts.Contact
	cursor   int
	offset   int
	height   int
	status   int
}

func NewModel(contacts []contacts.Contact) model {
	m := model{}
	m.contacts = contacts
	m.offset = 0
	m.cursor = 0
	m.height = 20
	m.status = StatusList
	return m
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) Home() {
	m.offset = 0
	m.cursor = 0
}

func (m *model) End() {
	m.cursor = len(m.contacts) - 1
	m.offset = m.cursor - m.height + 1
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *model) Up() {
	m.cursor = m.cursor - 1
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
}

func (m *model) Down() {
	m.cursor = m.cursor + 1
	if m.cursor > len(m.contacts)-1 {
		m.cursor = len(m.contacts) - 1
	}
	if m.cursor-m.height >= m.offset {
		m.offset = m.offset + 1
	}
}

func (m *model) PageDown() {
	m.offset = m.offset + m.height
	m.cursor = m.cursor + m.height

	if m.offset+m.height > len(m.contacts)-1 {
		m.offset = len(m.contacts) - m.height
		if m.offset < 0 {
			m.offset = 0
		}
	}

	if m.cursor > len(m.contacts)-1 {
		m.cursor = len(m.contacts) - 1
	}
}

func (m *model) PageUp() {
	m.offset = m.offset - m.height
	m.cursor = m.cursor - m.height

	if m.offset < 0 {
		m.offset = 0
	}

	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *model) Delete() {
	m.status = StatusConfirm
}

func (m model) UpdateList(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "pgdown":
			m.PageDown()
		case "pgup":
			m.PageUp()
		case "delete":
			m.Delete()
		}
	}
	return m, nil
}

func (m model) UpdateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "N", "n":
			m.status = StatusList
		case "Y", "y":
			m.status = StatusList
			db.Delete(context.Background(), m.contacts[m.cursor].ID)
			m.contacts, err = db.List(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if m.cursor > len(m.contacts)-1 {
				m.cursor = len(m.contacts) - 1
			}
		}
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.status {
	case StatusList:
		return m.UpdateList(msg)
	case StatusConfirm:
		return m.UpdateConfirm(msg)
	}
	return m, nil
}

func (m model) ViewList() string {
	var sb strings.Builder

	for i := m.offset; i < m.offset+m.height; i++ {
		if i >= len(m.contacts) {
			break
		}
		c := m.contacts[i]
		sb.WriteString(baseStyle.Reverse(i == m.cursor).Render(fmt.Sprintf("%-15s | %-15s | %-15s | %s", c.Lname, c.Fname, c.Phone, c.Email)) + "\n")
	}

	// details
	if len(m.contacts) > 0 {
		d := m.contacts[m.cursor]
		sb.WriteString(strings.Repeat("=", 80))
		sb.WriteString(fmt.Sprintf("\n%s %s, %s %s\n\n%s\n", d.Address, d.City, d.State, d.Zipcode, d.Notes))
	} else {
		sb.WriteString("No contacts")
	}

	return sb.String()
}

func (m model) ViewConfirm() string {
	c := m.contacts[m.cursor]
	style := lipgloss.NewStyle().
		Reverse(true).
		Border(lipgloss.RoundedBorder()).
		MarginLeft(3).
		Padding(0, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	v := m.ViewList()
	vLines := strings.Split(v, "\n")

	d := style.Render(fmt.Sprintf("Delete %s %s (Y/N)", c.Fname, c.Lname))
	dLines := strings.Split(d, "\n")

	for i, line := range dLines {
		vLines[i] = line
	}

	return strings.Join(vLines, "\n")
}

func (m model) View() string {
	var view string
	switch m.status {
	case StatusList:
		view = m.ViewList()
	case StatusConfirm:
		view = m.ViewConfirm()
	}
	return view
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

	db = contacts.New(conn)

	contacts, err := db.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	m := NewModel(contacts)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
