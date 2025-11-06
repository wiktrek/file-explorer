package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() model {
	return model{
		currentDir: config.DefaultPath,
		files:      loadFiles(config.DefaultPath),
		selected:   make(map[int]struct{}),
		keybinds:   config.keybinds,
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
		if msg.String() == "ctrl+q" || msg.String() == "q" {
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
					m.secondCursor = 1
				case "f2":
					m.viewState = Rename
					m.temp_string = m.files[m.cursor]
					m.secondCursor = len(m.temp_string)
				case "o":
					fileDir := m.currentDir + m.files[m.cursor]
					openFile(fileDir)
				case "ctrl+x", "m":
					m.viewState = Move
					m.temp_string = m.currentDir + m.files[m.cursor]
				case "ctrl+c":
					m.viewState = Copy
					m.temp_string = m.currentDir + m.files[m.cursor]
				case "ctrl+n":
					m.viewState = New
					m.temp_string = ""
				case "enter":
					pathToOpen := m.currentDir + m.files[m.cursor]
					v, err := IsDirectory(pathToOpen)
					if err != nil {
						fmt.Printf("%v", err)
					}
					if v {
						m = reloadDir(m, pathToOpen+"/")
					} else {
						fileDir := m.currentDir + m.files[m.cursor]
						openFile(fileDir)
					}

				case "esc":
					d := goUp(m.currentDir)
					m = reloadDir(m, d)
				}
			}
		case New:
			{
				switch msg.String() {
				case "left", "h":
					if m.secondCursor != 0 {
						m.secondCursor = 0
					}
				case "right", "l":
					if m.secondCursor != 1 {
						m.secondCursor = 1
					}
				case "backspace":
					m.temp_string = remove_at_index(m.temp_string, m.secondCursor)
					if m.secondCursor > 0 {
						m.secondCursor--
					}
				case "enter":
					createFile(m.currentDir + m.temp_string)
					m = reloadDir(m, "")
					m.viewState = Default
					m.View()
				default:
					if len(msg.String()) == 1 {
						m.temp_string = add_to_string(m.temp_string, msg.String()[0], m.secondCursor)
						m.secondCursor++
					}
				}
			}
		case ConfirmDelete:
			{

				switch msg.String() {
				case "left", "h":
					if m.secondCursor != 0 {
						m.secondCursor = 0
					}
				case "right", "l":
					if m.secondCursor != 1 {
						m.secondCursor = 1
					}

				case "enter":
					if m.secondCursor == 0 {
						deletePath(m.currentDir + m.files[m.cursor])

					}
					m = reloadDir(m, "")
					m.viewState = Default
					m.View()
				}

			}
		case Rename:
			{

				switch msg.String() {
				case "left", "h":
					if m.secondCursor > 0 {
						m.secondCursor--
					} else {
						m.secondCursor = len(m.temp_string)
					}
				case "right", "l":
					if m.secondCursor < len(m.temp_string) {
						m.secondCursor++
					} else {
						m.secondCursor = 0
					}
				case "up", "k":
					if m.cursor > 0 {
						m.cursor--
					} else {
						m.cursor = len(m.files) - 1
					}
					m.temp_string = m.files[m.cursor]
					if m.secondCursor > len(m.temp_string) {
						m.secondCursor = len(m.temp_string)
					}
				case "down", "j":
					if m.cursor < len(m.files)-1 {
						m.cursor++
					} else {
						m.cursor = 0
					}
					m.temp_string = m.files[m.cursor]
					if m.secondCursor > len(m.temp_string) {
						m.secondCursor = len(m.temp_string)
					}
				case "backspace":
					m.temp_string = remove_at_index(m.temp_string, m.secondCursor)
					if m.secondCursor > 0 {
						m.secondCursor--
					}
				case "enter":
					moveFile(m.currentDir+m.files[m.cursor], m.currentDir+m.temp_string)
					t := m.cursor
					reloadDir(m, "")
					m.viewState = Default
					m.View()
					m.secondCursor = 0
					m.cursor = t
					m.temp_string = ""
				default:
					if len(msg.String()) == 1 {
						m.temp_string = add_to_string(m.temp_string, msg.String()[0], m.secondCursor)
						m.secondCursor++
					}
				}
			}
		case Move:
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
				case "p":
					m.viewState = Default
				case "enter":
					pathToOpen := m.currentDir + m.files[m.cursor]
					v, err := IsDirectory(pathToOpen)
					if err != nil {
						fmt.Printf("%v", err)
					}
					if v {
						m = reloadDir(m, pathToOpen+"/")
					}
				case "ctrl+v":
					split := strings.Split(m.temp_string, "/")
					file_name := split[len(split)-1]
					err := moveFile(m.temp_string, m.currentDir+"/"+file_name)
					if err != nil {
						fmt.Println(err)
					} else {
						m.viewState = Default
					}
					m.temp_string = ""
					m.viewState = Default
					reloadDir(m, m.currentDir)
				case "esc":
					d := goUp(m.currentDir)
					m = reloadDir(m, d)
				}
			}
		case Copy:
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
				case "ctrl+v":
					split := strings.Split(m.temp_string, "/")
					file_name := split[len(split)-1]
					_, err := copyFile(m.temp_string, m.currentDir+file_name)
					if err != nil {
						fmt.Println(err)
					}
					m.temp_string = ""
					m.viewState = Default
					reloadDir(m, "")
				case "enter":
					pathToOpen := m.currentDir + m.files[m.cursor]
					v, err := IsDirectory(pathToOpen)
					if err != nil {
						fmt.Printf("%v", err)
					}
					if v {
						m = reloadDir(m, pathToOpen+"/")
					}
				case "esc":
					d := goUp(m.currentDir)
					m = reloadDir(m, d)
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
	case Move:
		return moveView(m)
	case Copy:
		return copyView(m)
	case New:
		return newView(m)
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
