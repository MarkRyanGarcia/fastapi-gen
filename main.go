package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markryangarcia/fastapi-gen/generator"
	"github.com/markryangarcia/fastapi-gen/tui"
)

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	cyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	muted  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	border = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	check  = lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true)
	errSty = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

func pipe() string { return border.Render("│") }

func main() {
	var initialName string

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "." {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println(errSty.Render("❌ Could not get current directory: " + err.Error()))
				os.Exit(1)
			}
			initialName = filepath.Base(cwd)
		} else {
			initialName = arg
		}
	}

	p := tea.NewProgram(tui.InitialModelWithName(initialName))
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println(errSty.Render("Error: " + err.Error()))
		os.Exit(1)
	}

	m := finalModel.(tui.Model)

	if m.Selected == "" || m.Quitting {
		fmt.Println(muted.Render("\n  Generation cancelled."))
		return
	}

	// Print vite-style summary
	fmt.Println(m.Summary())

	isSQL := strings.Contains(m.Selected, "SQL")
	isMongo := strings.Contains(m.Selected, "MongoDB")

	outDir := m.ProjectName
	if len(os.Args) > 1 && os.Args[1] == "." {
		outDir = "."
	}

	config := generator.ProjectConfig{
		ProjectName:       m.ProjectName,
		OutputDir:         outDir,
		Database:          m.Selected,
		IncludeSQLAlchemy: isSQL,
		IncludeMongoDB:    isMongo,
		AuthProvider:      m.AuthProvider,
		UseClerk:          m.AuthProvider == "Clerk",
		UseCognito:        m.AuthProvider == "AWS Cognito",
		UsePipenv:         m.UsePipenv,
		SetupVenv:         m.SetupVenv,
	}

	if err := generator.CreateProject(config); err != nil {
		fmt.Println(errSty.Render("❌ Failed to create project: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(check.Render("◇  ") + green.Render("Done! Project generated in ./"+outDir))
	fmt.Println(pipe())

	if m.SetupVenv {
		// activate + run
		fmt.Println(pipe() + "  " + cyan.Render("Starting dev server..."))
		fmt.Println(pipe())
		if err := generator.RunDevServer(outDir, m.UsePipenv); err != nil {
			fmt.Println(errSty.Render("❌ Failed to start dev server: " + err.Error()))
			os.Exit(1)
		}
	} else {
		// just print next steps
		fmt.Println(pipe() + "  " + cyan.Render("Next steps:"))
		if outDir != "." {
			fmt.Println(pipe() + "  " + muted.Render("cd "+outDir))
		}
		if m.UsePipenv {
			fmt.Println(pipe() + "  " + muted.Render("pipenv install"))
			fmt.Println(pipe() + "  " + muted.Render("pipenv shell"))
		} else {
			fmt.Println(pipe() + "  " + muted.Render("pip install -r requirements.txt"))
			fmt.Println(pipe() + "  " + muted.Render("source .venv/bin/activate"))
		}
		fmt.Println(pipe() + "  " + muted.Render("fastapi dev app"))
		fmt.Println(pipe())
	}
}
