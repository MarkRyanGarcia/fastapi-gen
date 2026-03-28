package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markryangarcia/fastapi-gen/generator"
	"github.com/markryangarcia/fastapi-gen/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var initialName string

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "." {
			// Use current directory name as project name, scaffold in-place
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("❌ Could not get current directory: %v\n", err)
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
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	m := finalModel.(tui.Model)

	if m.Selected != "" && !m.Quitting {
		isSQL := strings.Contains(m.Selected, "SQL")
		isMongo := strings.Contains(m.Selected, "MongoDB")

		// Determine output directory
		outDir := m.ProjectName
		if len(os.Args) > 1 && os.Args[1] == "." {
			outDir = "."
		}

		fmt.Printf("\n🚀 Creating project '%s'...\n", m.ProjectName)

		config := generator.ProjectConfig{
			ProjectName:       m.ProjectName,
			OutputDir:         outDir,
			Database:          m.Selected,
			IncludeSQLAlchemy: isSQL,
			IncludeMongoDB:    isMongo,
			UsePipenv:         m.UsePipenv,
		}

		err := generator.CreateProject(config)
		if err != nil {
			fmt.Printf("❌ Failed to create project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Success! Project generated in ./%s\n", outDir)
	} else {
		fmt.Println("\nGeneration cancelled.")
	}
}
