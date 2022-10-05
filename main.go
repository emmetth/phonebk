package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emmetth/phonebk/contacts"
	_ "github.com/mattn/go-sqlite3"
)

func NewModel() (*model, error) {
	return &model{}, nil
}

type model struct {
	contacts []contacts.Contact
	cursor   int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.contacts)-1 {
				m.cursor++
			}

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := "" //user string builder

	// Iterate over our choices
	for i, contact := range m.contacts {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, contact.Fname)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
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

	m := model{contacts: contacts}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
