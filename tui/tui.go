package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── styles ────────────────────────────────────────────────────────────────────

var (
	cyan      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	green     = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	muted     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	bold      = lipgloss.NewStyle().Bold(true)
	cursorSty = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	checkSty  = lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true)
	labelSty  = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	valueSty  = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	borderSty = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// ── state machine ─────────────────────────────────────────────────────────────

type sessionState int

const (
	stateInputName sessionState = iota
	stateSelectDB
	stateSelectAuth
	stateSelectPipenv
	stateSelectVenv
	stateDone
)

// ── model ─────────────────────────────────────────────────────────────────────

type Model struct {
	State        sessionState
	TextInput    textinput.Model
	ProjectName  string
	Cursor       int
	Choices      []string
	Selected     string
	AuthProvider string
	UsePipenv    bool
	SetupVenv    bool
	Quitting     bool
}

func InitialModel() Model {
	return InitialModelWithName("")
}

func InitialModelWithName(name string) Model {
	ti := textinput.New()
	ti.Placeholder = "my-fastapi-app"
	ti.CharLimit = 32
	ti.Width = 24
	ti.PromptStyle = cyan
	ti.TextStyle = valueSty

	state := stateInputName
	if name != "" {
		state = stateSelectDB
	} else {
		ti.Focus()
	}

	return Model{
		State:       state,
		TextInput:   ti,
		ProjectName: name,
		Choices:     []string{"PostgreSQL (SQLAlchemy)", "MongoDB (PyMongo)"},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// ── update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			switch m.State {
			case stateInputName:
				m.ProjectName = m.TextInput.Value()
				if m.ProjectName == "" {
					m.ProjectName = "my-fastapi-app"
				}
				m.State = stateSelectDB
				m.Cursor = 0
			case stateSelectDB:
				m.Selected = m.Choices[m.Cursor]
				m.State = stateSelectAuth
				m.Cursor = 0
			case stateSelectAuth:
				authChoices := []string{"None", "Clerk", "AWS Cognito"}
				m.AuthProvider = authChoices[m.Cursor]
				m.State = stateSelectPipenv
				m.Cursor = 0
			case stateSelectPipenv:
				m.UsePipenv = m.Cursor == 0
				m.State = stateSelectVenv
				m.Cursor = 0
			case stateSelectVenv:
				m.SetupVenv = m.Cursor == 0
				m.State = stateDone
				return m, tea.Quit
			}
			return m, nil

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			switch m.State {
			case stateSelectDB:
				if m.Cursor < len(m.Choices)-1 {
					m.Cursor++
				}
			case stateSelectAuth:
				if m.Cursor < 2 {
					m.Cursor++
				}
			case stateSelectPipenv, stateSelectVenv:
				if m.Cursor < 1 {
					m.Cursor++
				}
			}
		}
	}

	if m.State == stateInputName {
		m.TextInput, cmd = m.TextInput.Update(msg)
	}

	return m, cmd
}

// ── helpers ───────────────────────────────────────────────────────────────────

func pkgManagerLabel(usePipenv bool) string {
	if usePipenv {
		return "pipenv"
	}
	return "pip"
}

func renderOptions(choices []string, cursor int) string {
	var sb strings.Builder
	for i, choice := range choices {
		if i == cursor {
			sb.WriteString(cursorSty.Render("❯ ") + bold.Render(choice) + "\n")
		} else {
			sb.WriteString(muted.Render("  "+choice) + "\n")
		}
	}
	return sb.String()
}

func pipe() string { return borderSty.Render("│") }

func summaryRow(label, value string) string {
	return fmt.Sprintf("%s  %s %s\n", pipe(), labelSty.Render(label), valueSty.Render(value))
}

// ── view ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.Quitting {
		return cyan.Render("  Cancelled.\n")
	}

	hint := muted.Render("  ↑/↓ or j/k to move · enter to select\n")

	switch m.State {
	case stateInputName:
		return fmt.Sprintf(
			"%s\n%s\n\n  %s\n\n%s",
			pipe(),
			pipe()+"  "+cyan.Render("Project name:"),
			m.TextInput.View(),
			muted.Render("  enter to confirm\n"),
		)

	case stateSelectDB:
		return fmt.Sprintf(
			"%s\n%s\n\n%s\n%s",
			pipe(),
			pipe()+"  "+cyan.Render("Select a database:"),
			renderOptions(m.Choices, m.Cursor),
			hint,
		)

	case stateSelectAuth:
		return fmt.Sprintf(
			"%s\n%s\n\n%s\n%s",
			pipe(),
			pipe()+"  "+cyan.Render("Select an auth provider:"),
			renderOptions([]string{"None", "Clerk", "AWS Cognito"}, m.Cursor),
			hint,
		)

	case stateSelectPipenv:
		return fmt.Sprintf(
			"%s\n%s\n\n%s\n%s",
			pipe(),
			pipe()+"  "+cyan.Render("Package manager:"),
			renderOptions([]string{"Pipenv", "requirements.txt"}, m.Cursor),
			hint,
		)

	case stateSelectVenv:
		return fmt.Sprintf(
			"%s\n%s\n\n%s\n%s",
			pipe(),
			pipe()+"  "+cyan.Render("Install with "+pkgManagerLabel(m.UsePipenv)+" and start now?"),
			renderOptions([]string{"Yes", "No"}, m.Cursor),
			hint,
		)
	}

	return ""
}

// ── summary (called from main after program exits) ────────────────────────────

func (m Model) Summary() string {
	if m.Quitting || m.Selected == "" {
		return ""
	}

	pkgManager := "requirements.txt"
	if m.UsePipenv {
		pkgManager = "Pipenv"
	}
	installNow := "No"
	if m.SetupVenv {
		installNow = "Yes"
	}

	var sb strings.Builder
	sb.WriteString(pipe() + "\n")
	sb.WriteString(summaryRow("Project:       ", m.ProjectName))
	sb.WriteString(summaryRow("Database:      ", m.Selected))
	sb.WriteString(summaryRow("Auth:          ", m.AuthProvider))
	sb.WriteString(summaryRow("Pkg manager:   ", pkgManager))
	sb.WriteString(summaryRow("Install & start: ", installNow))
	sb.WriteString(pipe() + "\n")
	sb.WriteString(checkSty.Render("◇  ") + green.Render("Scaffolding project in ./"+m.ProjectName+"...") + "\n")

	return sb.String()
}
