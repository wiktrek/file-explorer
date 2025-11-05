package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() model {
	return model{
		currentDir: defaultDir,
		files:      loadFiles(defaultDir),
		selected:   make(map[int]struct{}),
		viewState:  Default,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// I don't want to add quitting to every View
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		switch m.viewState {
		case Default:
			{
				switch msg.String() {
				case "up", "k":
					if m.cursor > 0 {
						m.cursor--
					} else {
						m.cursor = len(m.files) - 1
					}

				case "down", "j":
					if m.cursor < len(m.files)-1 {
						m.cursor++
					} else {
						m.cursor = 0
					}
				case "d":
					m.viewState = ConfirmDelete
					m.confirmCursor = 1
				case "f2":
					m.viewState = Rename
					m.temp_string = m.files[m.cursor]
					m.confirmCursor = len(m.temp_string)
				case "o", " ":
					fileDir := m.currentDir + m.files[m.cursor]
					openFile(fileDir)
				case "enter":
					pathToOpen := m.currentDir + m.files[m.cursor]
					v, err := IsDirectory(pathToOpen)
					if err != nil {
						fmt.Printf("%v", err)
					}
					if v {
						m = reloadDir(m, pathToOpen+"/")
					} else {
						// idk what to do for enter
					}

				case "esc":
					d := goUp(m.currentDir)
					m = reloadDir(m, d)
				}
			}
		case ConfirmDelete:
			{

				switch msg.String() {
				case "left", "h":
					if m.confirmCursor != 0 {
						m.confirmCursor = 0
					}
				case "right", "l":
					if m.confirmCursor != 1 {
						m.confirmCursor = 1
					}
				}
			}
		case Rename:
			{

				switch msg.String() {
				case "left", "h":
					if m.confirmCursor > 0 {
						m.confirmCursor--
					} else {
						m.confirmCursor = len(m.temp_string)
					}
				case "right", "l":
					if m.confirmCursor < len(m.temp_string) {
						m.confirmCursor++
					} else {
						m.confirmCursor = 0
					}
				}

			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.viewState {
	case Default:
		return defaultView(m)
	case ConfirmDelete:
		return confirmDeleteView(m)
	case Rename:
		return renameView(m)
	}
	return "Loading..."
}
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
