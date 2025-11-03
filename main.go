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

type model struct {
	files      []string
	cursor     int
	selected   map[int]struct{}
	currentDir string
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

	fmt.Println(split)
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
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

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
				m.cursor = 0
				m.currentDir = pathToOpen + "/"
				m.files = loadFiles(m.currentDir)

			} else {
				// idk what to do for enter
			}

		case "esc":
			d := goUp(m.currentDir)
			fmt.Print(d)
			m.currentDir = d
			m.files = loadFiles(d)
			m.View()
		}

	}

	return m, nil
}
func (m model) View() string {
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
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
