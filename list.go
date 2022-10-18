package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emmetth/phonebk/contacts"
	_ "github.com/mattn/go-sqlite3"
)

type ListModel struct {
	contacts []contacts.Contact
	cursor   int
	offset   int
	height   int
	confirm  bool
}

func NewListModel(contacts []contacts.Contact) ListModel {
	m := ListModel{}
	m.contacts = contacts
	m.offset = 0
	m.cursor = 0
	m.height = 20
	m.confirm = false
	return m
}

func (m ListModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *ListModel) Home() {
	m.offset = 0
	m.cursor = 0
}

func (m *ListModel) End() {
	m.cursor = len(m.contacts) - 1
	m.offset = m.cursor - m.height + 1
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *ListModel) Up() {
	m.cursor = m.cursor - 1
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
}

func (m *ListModel) Down() {
	m.cursor = m.cursor + 1
	if m.cursor > len(m.contacts)-1 {
		m.cursor = len(m.contacts) - 1
	}
	if m.cursor-m.height >= m.offset {
		m.offset = m.offset + 1
	}
}

func (m *ListModel) PageDown() {
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

func (m *ListModel) PageUp() {
	m.offset = m.offset - m.height
	m.cursor = m.cursor - m.height

	if m.offset < 0 {
		m.offset = 0
	}

	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m ListModel) UpdateList(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.confirm = true
		case "enter":
			return m, EditCmd(m.contacts[m.cursor])
		case "insert", "alt+[2~":
			return m, EditCmd(contacts.Contact{})
		}

		if msg.Type == tea.KeyRunes {
			key := strings.ToLower(msg.String())
			if len(key) == 1 && key >= "a" && key <= "z" {
				newCursor := len(m.contacts) - 1
				for i, contact := range m.contacts {
					if strings.Compare(key, strings.ToLower(contact.Fname)) <= 0 {
						newCursor = i
						break
					}
				}
				m.cursor = newCursor
			}
		}
	}
	return m, nil
}

func (m ListModel) UpdateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "N", "n":
			m.confirm = false
		case "Y", "y":
			m.confirm = false
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

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height - 5
		return m, nil
	}
	if m.confirm {
		return m.UpdateConfirm(msg)
	}
	return m.UpdateList(msg)
}

func (m ListModel) ViewList() string {
	var sb strings.Builder

	for i := m.offset; i < m.offset+m.height; i++ {
		if i >= len(m.contacts) {
			break
		}
		c := m.contacts[i]
		sb.WriteString(baseStyle.Reverse(i == m.cursor).Render(fmt.Sprintf("%-25s | %-12s | %s", c.Fname+" "+c.Lname, c.Phone, c.Email)) + "\n")
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

func (m ListModel) ViewConfirm() string {
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

func (m ListModel) View() string {
	var view string
	if m.confirm {
		view = m.ViewConfirm()
	} else {
		view = m.ViewList()
	}
	return view
}
