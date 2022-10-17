package main

import (
	"context"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emmetth/phonebk/contacts"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type EditModel struct {
	focus     int
	contactId int64
	inputs    []textinput.Model
}

func (m EditModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func NewEditModel(contact contacts.Contact) EditModel {
	m := EditModel{}
	m.focus = 0
	m.contactId = contact.ID
	m.inputs = make([]textinput.Model, 10)

	for i := range m.inputs {
		t := textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 40
		t.Width = 40

		switch i {
		case 0:
			t.Placeholder = "First Name"
			t.Focus()
			t.SetValue(contact.Fname)
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
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
		m.inputs[i] = t
	}

	return m
}

func (m EditModel) contact() contacts.Contact {
	return contacts.Contact{
		Fname:    m.inputs[0].Value(),
		Lname:    m.inputs[1].Value(),
		Phone:    m.inputs[2].Value(),
		Email:    m.inputs[3].Value(),
		Address:  m.inputs[4].Value(),
		City:     m.inputs[5].Value(),
		State:    m.inputs[6].Value(),
		Zipcode:  m.inputs[7].Value(),
		Birthday: m.inputs[8].Value(),
		Notes:    m.inputs[9].Value(),
		ID:       m.contactId}
}

func (m EditModel) View() string {
	var sb strings.Builder

	for _, input := range m.inputs {
		sb.WriteString(input.View() + "\n")
	}

	return sb.String()
}

func (m *EditModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *EditModel) setFocus() {
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focus {
			// Set focused state
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}
}

func (m *EditModel) down() {
	m.focus++
	if m.focus >= len(m.inputs) {
		m.focus = len(m.inputs) - 1
	}
	m.setFocus()
}

func (m *EditModel) up() {
	m.focus--
	if m.focus < 0 {
		m.focus = 0
	}
	m.setFocus()
}

func (m EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			c := m.contact()
			if c.ID == 0 {
				params := contacts.AddParams{Fname: c.Fname, Lname: c.Lname, Phone: c.Phone, Email: c.Email, Address: c.Address, State: c.State, Zipcode: c.Zipcode, Birthday: c.Birthday, Notes: c.Notes}
				err := db.Add(context.Background(), params)
				if err != nil {
					log.Panic(err)
				}

			} else {
				params := contacts.UpdateParams{ID: c.ID, Fname: c.Fname, Lname: c.Lname, Phone: c.Phone, Email: c.Email, Address: c.Address, State: c.State, Zipcode: c.Zipcode, Birthday: c.Birthday, Notes: c.Notes}
				err := db.Update(context.Background(), params)
				if err != nil {
					log.Panic(err)
				}
			}
			return m, BackCmd()
		case "up", "shift+tab":
			m.up()
		case "enter", "down", "tab":
			m.down()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}
