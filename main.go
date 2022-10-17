package main

import (
	"context"
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emmetth/phonebk/contacts"
	_ "github.com/mattn/go-sqlite3"
)

var baseStyle = lipgloss.NewStyle()
var db *contacts.Queries

const (
	StateList = iota
	StateEdit
)

type EditMsg struct {
	contact contacts.Contact
}

type BackMsg struct {
	contact contacts.Contact
}

type MainModel struct {
	state int
	list  ListModel
	edit  EditModel
}

func EditCmd(contact contacts.Contact) tea.Cmd {
	return func() tea.Msg {
		return EditMsg{contact}
	}
}

func BackCmd(contact contacts.Contact) tea.Cmd {
	return func() tea.Msg {
		return BackMsg{contact}
	}
}

func NewMainModel(contacts []contacts.Contact) MainModel {
	m := MainModel{}
	m.list = NewListModel(contacts)
	return m
}

func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case EditMsg:
		m.state = StateEdit
		m.edit = NewEditModel(msg.(EditMsg).contact)
	case BackMsg:
		m.state = StateList
		m.list.contacts[m.list.cursor] = msg.(BackMsg).contact
	}

	switch m.state {
	case StateList:
		newList, newCmd := m.list.Update(msg)
		m.list = newList.(ListModel)
		return m, newCmd
	case StateEdit:
		newEdit, newCmd := m.edit.Update(msg)
		m.edit = newEdit.(EditModel)
		return m, newCmd
	}

	return m, nil
}

func (m MainModel) View() string {
	var view string
	switch m.state {
	case StateList:
		view = m.list.View()
	case StateEdit:
		view = m.edit.View()
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

	m := NewMainModel(contacts)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
