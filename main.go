package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emmetth/phonebk/contacts"
	_ "github.com/mattn/go-sqlite3"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewModel() (*model, error) {
	return &model{}, nil
}

type model struct {
	contacts []contacts.Contact
	table    table.Model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + baseStyle.Render(m.Details()) + "\n"
}

func (m model) Details() string {
	c := m.contacts[m.table.Cursor()]
	return fmt.Sprintf("%s %s, %s %s\n\n%s\n", c.Address, c.City, c.State, c.Zipcode, c.Notes)
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

	columns := []table.Column{
		{Title: "First", Width: 15},
		{Title: "Last", Width: 15},
		{Title: "Phone", Width: 15},
		{Title: "Email", Width: 40},
	}

	var rows []table.Row

	for _, contact := range contacts {
		rows = append(rows, table.Row{contact.Fname, contact.Lname, contact.Phone, contact.Email})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	m := model{contacts, t}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
