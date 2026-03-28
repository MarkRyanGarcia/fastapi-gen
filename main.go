package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/markryangarcia/fastapi-gen/generator"
	"github.com/markryangarcia/fastapi-gen/tui" 
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	m := finalModel.(tui.Model)

	if m.Selected != "" && !m.Quitting {
		fmt.Printf("\n🚀 Creating project '%s'...\n", m.ProjectName)
		
		isSQL := strings.Contains(m.Selected, "SQL") || strings.Contains(m.Selected, "SQLite")
		isMongo := strings.Contains(m.Selected, "MongoDB")

		config := generator.ProjectConfig{
			ProjectName:       m.ProjectName,
			Database:          m.Selected,
			IncludeSQLAlchemy: isSQL,
			IncludeMongoDB:    isMongo,
		}

		err := generator.CreateProject(config)
		if err != nil {
			fmt.Printf("❌ Failed to create project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Success! Project generated in ./%s\n", m.ProjectName)
	} else {
		fmt.Println("\nGeneration cancelled.")
	}
}