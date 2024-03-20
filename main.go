package main


import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const spinner = "Loading..."
const quitToolTip = "(ctrl+c to quit)"

type Styles struct {
	BorderColor    lipgloss.Color
	InputField     lipgloss.Style
	headerRowStyle lipgloss.Style
	baseRowStyle   lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	s.baseRowStyle = lipgloss.NewStyle().Padding(0, 1)
	s.headerRowStyle = s.baseRowStyle.Copy().Foreground(lipgloss.Color("252")).Bold(true)
	return s
}

type model struct {
	index     int
	questions []Question
	width     int
	height    int
	styles    *Styles
	done      bool
}

type Question struct {
	question string
	answer   string
	input    Input
}

func NewQuestion(question string) Question {
	return Question{question: question}
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
		questions: questions,
		styles:    styles,
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
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyShiftTab.String():
			m.Previous()
		case tea.KeyTab.String():
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.answer = current.input.Value()
			m.Next()

			return m, current.input.Blur
		}

	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {

	if m.width == 0 {
		return spinner
	}

	if m.done == true {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				m.Result().String(),
				quitToolTip,
			),
		)
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
			quitToolTip,
		),
	)

}

func (m *model) Result() *table.Table {
	var rows [][]string

	headers := []string{"#", "Q", "A"}
	for i, q := range m.questions {
		log.Println(i)
		rows = append(rows, []string{
			fmt.Sprintf("%d",i),
			q.question,
			q.answer,
		})
	}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(80).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return m.styles.headerRowStyle
			}

			even := row%2 == 0

			if even {
				return m.styles.baseRowStyle.Copy().Foreground(lipgloss.Color("245"))
			}
			return m.styles.baseRowStyle.Copy().Foreground(lipgloss.Color("252"))
		})
	return t
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func (m *model) Previous() {
	if m.index > 0 {
		m.index--
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
