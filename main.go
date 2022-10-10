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

type BackMsg struct{}

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

func BackCmd() tea.Msg {
	return BackMsg{}
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
	log.Printf("main msg: %+v\n", msg)
	switch msg.(type) {
	case EditMsg:
		m.state = StateEdit
		log.Println(msg)
	case BackMsg:
		m.state = StateList
	}

	switch m.state {
	case StateList:
		return m.list.Update(msg)
	case StateEdit:
		return m.edit.Update(msg)
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
