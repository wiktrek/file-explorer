package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewState string

const (
	Default       ViewState = "default"
	ConfirmDelete ViewState = "confirm"
)

type model struct {
	files         []string
	cursor        int
	selected      map[int]struct{}
	currentDir    string
	confirmCursor int
	viewState     ViewState
}

func reloadDir(m model, path string) {
	m.cursor = 0
	if path == "" {
		m.currentDir = path
	}
	m.files = loadFiles(m.currentDir)
	m.View()
}
func deleteFile(path string) {
	os.Remove(path)

}
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
func openFile(f string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, f)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
func loadFiles(dir string) []string {
	pattern := "*"
	matches, err := filepath.Glob(dir + pattern)
	if err != nil {
		log.Fatal(err)
	}
	var files []string
	for _, match := range matches {
		split := strings.Split(match, "/")
		files = append(files, split[len(split)-1])
	}
	return files
}
func goUp(dir string) string {
	if dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	var split = strings.Split(dir, "/")
	result := strings.Join(split[:len(split)-1], "/")

	if len(result) <= 1 {
		return "/"
	}
	return result + "/"
}

var defaultDir = "/home/wiktor/projects/golang/file-explorer/"

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
		if m.viewState == Default {
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.files)-1 {
					m.cursor++
				}
			case "d":
				m.viewState = ConfirmDelete
				m.confirmCursor = 1
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
					reloadDir(m, pathToOpen+"/")
				} else {
					// idk what to do for enter
				}

			case "esc":
				d := goUp(m.currentDir)
				reloadDir(m, d)
			}
		} else {

			switch msg.String() {
			case "left", "h":
				if m.confirmCursor != 0 {
					m.confirmCursor = 0
				}
			case "right", "l":
				if m.confirmCursor != 1 {
					m.confirmCursor = 1
				}
			case "enter":
				if m.confirmCursor == 0 {
					deleteFile(m.currentDir + m.files[m.cursor])

				}
				m.viewState = Default
				m.View()
			}
		}
	}
	return m, nil
}
func defaultView(m model) string {
	s := "Select File: \n\n"
	s += m.currentDir + "\n"

	for i, file := range m.files {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, file)
	}
	s += "\nPress q to quit.\nPress o to open selected file\n"
	return s
}
func confirmDeleteView(m model) string {
	s := "Are you really sure you want to delete"
	s += "\n" + m.currentDir + m.files[m.cursor] + "\n"
	if m.confirmCursor == 0 {
		s += ">"
	}
	s += "YES"
	s += "		"
	if m.confirmCursor == 1 {
		s += ">"
	}
	s += "NO"
	return s
}
func (m model) View() string {
	switch m.viewState {
	case Default:
		return defaultView(m)
	case ConfirmDelete:
		return confirmDeleteView(m)
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
