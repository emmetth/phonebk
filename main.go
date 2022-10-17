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

func BackCmd() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}

func NewMainModel() MainModel {
	m := MainModel{}
	m.Load()
	return m
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case EditMsg:
		m.state = StateEdit
		m.edit = NewEditModel(msg.(EditMsg).contact)
	case BackMsg:
		m.state = StateList
		m.Load()
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

func (m *MainModel) Load() {
	contacts, err := db.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	m.list = NewListModel(contacts)
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

	m := NewMainModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
