package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emmetth/phonebk/contacts"
)

type EditModel struct {
	focus  int
	inputs []textinput.Model
}

func (m EditModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func NewEditModel(contact contacts.Contact) EditModel {
	m := EditModel{}
	m.focus = 0
	m.inputs = make([]textinput.Model, 10)

	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = 40
		t.Width = 40

		switch i {
		case 0:
			t.Placeholder = "First Name"
			t.Focus()
			t.SetValue(contact.Fname)
		case 1:
			t.Placeholder = "Last Name"
			t.SetValue(contact.Lname)
		case 2:
			t.Placeholder = "Phone"
			t.SetValue(contact.Phone)
		case 3:
			t.Placeholder = "Email"
			t.SetValue(contact.Email)
		case 4:
			t.Placeholder = "Address"
			t.SetValue(contact.Address)
		case 5:
			t.Placeholder = "City"
			t.SetValue(contact.City)
		case 6:
			t.Placeholder = "State"
			t.SetValue(contact.State)
		case 7:
			t.Placeholder = "Zipcode"
			t.SetValue(contact.Zipcode)
		case 8:
			t.Placeholder = "mm/dd/yyyy"
			t.SetValue(contact.Birthday)
		case 9:
			t.Placeholder = "Notes"
			t.SetValue(contact.Notes)
		}
	}

	return m
}

func (m EditModel) View() string {
	var sb strings.Builder

	for i := range m.inputs {
		sb.WriteString(m.inputs[i].View() + "\n")
	}

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
