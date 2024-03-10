package main

// TODO: 30mn44s

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct {
	index       int
	questions   []Question
	width       int
	height      int
	styles      *Styles
}

type Question struct {
	question string
	answer string
	input Input
}


func NewQuestion(question string) Question {
	return Question{ question: question}
}

func NewShortQuestion(question string) Question {
	q := Question{question: question}
	q.input = NewShortAnswerField()
	return q	
}

func NewLongQuestion(question string) Question {
	q := Question{question: question}
	q.input = NewLongAnswerField()
	return q
}

func New(questions []Question) *model {
	styles := DefaultStyles()

	return &model{
		questions:   questions,
		styles:      styles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			current.answer = current.input.Value()
			current.input.Blur()

			log.Printf("question: %s, answer: %s", current.question, current.answer)
			m.Next()

			return m, nil
		}
		

	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	current := &m.questions[m.index]
	current.input.Focus()
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			current.question,
			m.styles.InputField.Render(current.input.View()),
		),
	)

}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {

	questions := []Question{
		NewShortQuestion("What is your name?"),
		NewShortQuestion("What is your favourite editor?"),
		NewLongQuestion("What is your favour quote?"),
	}

	m := New(questions)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
