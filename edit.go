package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emmetth/phonebk/contacts"
)

type EditModel struct {
	contact contacts.Contact
}

func (m EditModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func NewEditModel(contact contacts.Contact) EditModel {
	return EditModel{contact}
}

func (m EditModel) View() string {
	var sb strings.Builder
	sb.WriteString("edit screen")
	sb.WriteString("First Name: " + m.contact.Fname)
	return sb.String()
}

func (m EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// send back cmd
			return m, BackCmd()
		}
	}
	return m, nil
}
