package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Update(msg tea.Msg) (Input, tea.Cmd)
	Blur() tea.Msg
	View() string
	Focus() tea.Cmd
}

// ShortAnswer

type ShortAnswerField struct {
	textinput textinput.Model
}

func NewShortAnswerField() *ShortAnswerField {
	ti := textinput.New()
	ti.Placeholder = "Your answer here"
	return &ShortAnswerField{
		textinput: ti,
	}
}

func (sa *ShortAnswerField) Value() string {
	return sa.textinput.Value()
}

func (sa *ShortAnswerField) Blur() tea.Msg {
	return sa.textinput.Blur
}

func (sa *ShortAnswerField) Focus() tea.Cmd {
	return sa.textinput.Focus()
}

func (sa *ShortAnswerField) View() string {
	return sa.textinput.View()
}

func (sa *ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	sa.textinput, cmd = sa.textinput.Update(msg)

	return sa, cmd
}

// LongAnswer

type LongAnswerField struct {
	textarea textarea.Model
}

func NewLongAnswerField() *LongAnswerField {
	ta := textarea.New()
	ta.Placeholder = "Your answer here"
	return &LongAnswerField{
		textarea: ta,
	}
}

func (la *LongAnswerField) Value() string {
	return la.textarea.Value()
}

func (la *LongAnswerField) Blur() tea.Msg {
	return la.textarea.Blur
}

func (la *LongAnswerField) Focus() tea.Cmd {
	return la.textarea.Focus()
}

func (sa *LongAnswerField) View() string {
	return sa.textarea.View()
}

func (la *LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	la.textarea, cmd = la.textarea.Update(msg)
	return la, cmd
}
