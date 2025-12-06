package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialmodel() model {
	return model{
		currentDir: config.defaultPath,
		files:      loadFiles(config.defaultPath),
		selected:   make(map[int]struct{}),
		config:     config,
		searching:  false,
		viewState:  Default,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("File explorer")
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// I don't want to add quitting to every View
		if msg.String() == "ctrl+q" || msg.String() == "q" || msg.String() == "ctrl+w" {
			return m, tea.Quit
		}
		switch m.viewState {
		case Default:
			{
				if m.searching {
					switch msg.String() {
					case "esc":
						m.searching = false
						m.search = ""
					case "enter":
						m.searching = false
					case "up":
						if m.cursor > 0 {
							m.cursor--
						} else {
							m.cursor = len(m.files) - 1
						}
					case "down":
						if m.cursor < len(m.files)-1 {
							m.cursor++
						} else {
							m.cursor = 0
						}
					case "left":
						if m.secondCursor != 0 {
							m.secondCursor = 0
						}
					case "right":
						if m.secondCursor != 1 {
							m.secondCursor = 1
						}
					case "backspace":
						m.search = remove_at_index(m.search, m.secondCursor)
						if m.secondCursor > 0 {
							m.secondCursor--
						}
						m = reloadDir(m, m.currentDir)
					default:
						if len(msg.String()) == 1 {
							m.search = add_to_string(m.search, msg.String()[0], m.secondCursor)
							m.secondCursor++
							m = reloadDir(m, m.currentDir)
						}
					}
				} else {
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
					case "d", "delete":
						m.viewState = ConfirmDelete
						m.secondCursor = 1
					case "f2":
						m.viewState = Rename
						m.tempString = m.files[m.cursor].path
						m.secondCursor = len(m.tempString)
					case "o":
						fileDir := m.currentDir + m.files[m.cursor].path
						openFile(fileDir)
					case "ctrl+x", "m":
						m.viewState = Move
						m.tempString = m.currentDir + m.files[m.cursor].path
					case "ctrl+c":
						m.viewState = Copy
						m.tempString = m.currentDir + m.files[m.cursor].path
					case "ctrl+n":
						m.viewState = New
						m.tempString = ""
					case "enter":
						pathToOpen := m.currentDir + m.files[m.cursor].path
						v, err := IsDirectory(pathToOpen)
						if err != nil {
							fmt.Printf("%v", err)
						}
						if v {
							m = reloadDir(m, pathToOpen+"/")
						} else {
							fileDir := m.currentDir + m.files[m.cursor].path
							openFile(fileDir)
						}
					case "/":
						m.searching = true
					case "esc":
						m.search = ""
						d := goUp(m.currentDir)
						m = reloadDir(m, d)
					}
				}
			}
		case New:
			{
				switch msg.String() {
				case "left":
					if m.secondCursor != 0 {
						m.secondCursor = 0
					}
				case "right":
					if m.secondCursor != 1 {
						m.secondCursor = 1
					}
				case "backspace":
					m.tempString = remove_at_index(m.tempString, m.secondCursor)
					if m.secondCursor > 0 {
						m.secondCursor--
					}
				case "enter":
					createFile(m.currentDir + m.tempString)
					t := m.cursor
					m.secondCursor = 0
					m.cursor = t
					m.tempString = ""
					m = reloadDir(m, "")
					m.viewState = Default
					m.View()
				default:
					if len(msg.String()) == 1 {
						m.tempString = add_to_string(m.tempString, msg.String()[0], m.secondCursor)
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
						deletePath(m.currentDir + m.files[m.cursor].path)

					}
					m = reloadDir(m, "")
					m.viewState = Default
					m.View()
				}

			}
		case Rename:
			{

				switch msg.String() {
				case "left":
					if m.secondCursor > 0 {
						m.secondCursor--
					} else {
						m.secondCursor = len(m.tempString)
					}
				case "right":
					if m.secondCursor < len(m.tempString) {
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
					m.tempString = m.files[m.cursor].path
					if m.secondCursor > len(m.tempString) {
						m.secondCursor = len(m.tempString)
					}
				case "down", "j":
					if m.cursor < len(m.files)-1 {
						m.cursor++
					} else {
						m.cursor = 0
					}
					m.tempString = m.files[m.cursor].path
					if m.secondCursor > len(m.tempString) {
						m.secondCursor = len(m.tempString)
					}
				case "backspace":
					m.tempString = remove_at_index(m.tempString, m.secondCursor)
					if m.secondCursor > 0 {
						m.secondCursor--
					}
				case "enter":
					moveFile(m.currentDir+m.files[m.cursor].path, m.currentDir+m.tempString)
					t := m.cursor
					m = reloadDir(m, "")
					m.viewState = Default
					m.secondCursor = 0
					m.cursor = t
					m.tempString = ""
					m.View()
				default:
					if len(msg.String()) == 1 {
						m.tempString = add_to_string(m.tempString, msg.String()[0], m.secondCursor)
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
					pathToOpen := m.currentDir + m.files[m.cursor].path
					v, err := IsDirectory(pathToOpen)
					if err != nil {
						fmt.Printf("%v", err)
					}
					if v {
						m = reloadDir(m, pathToOpen+"/")
					}
				case "ctrl+v":
					split := strings.Split(m.tempString, "/")
					file_name := split[len(split)-1]
					err := moveFile(m.tempString, m.currentDir+"/"+file_name)
					if err != nil {
						fmt.Println(err)
					} else {
						m.viewState = Default
					}
					m.tempString = ""
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
					split := strings.Split(m.tempString, "/")
					file_name := split[len(split)-1]
					_, err := copyFile(m.tempString, m.currentDir+file_name)
					if err != nil {
						fmt.Println(err)
					}
					m.tempString = ""
					m.viewState = Default
					reloadDir(m, "")
				case "enter":
					pathToOpen := m.currentDir + m.files[m.cursor].path
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
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	var views []string

	switch m.viewState {
	case Default:
		views = append(views, defaultView(m))
	case ConfirmDelete:
		views = append(views, confirmDeleteView(m))
	case Rename:
		views = append(views, renameView(m))
	case Move:
		views = append(views, moveView(m))
	case Copy:
		views = append(views, copyView(m))
	case New:
		views = append(views, newView(m))
	}
	views = append(views, previewView(m))
	if len(views) != 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + showBinds(m.config.keybinds)
	}
	return "Loading..."
}
func main() {
	p := tea.NewProgram(initialmodel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
