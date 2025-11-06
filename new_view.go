package main

import "fmt"

func newView(m model) string {
	s := "Create new file: "
	s += m.currentDir + "\n"

	for _, file := range m.files {
		s += fmt.Sprintf("%s %s\n", file.fileType, file.path)
	}
	if m.secondCursor == len(m.temp_string) {
		s += m.temp_string
		s += "|"
	} else {
		for j := range m.temp_string {
			if j == m.secondCursor {
				s += fmt.Sprintf("|%c", m.temp_string[j])
			} else {
				s += fmt.Sprintf("%c", m.temp_string[j])
			}
		}
	}
	s += "\n"
	s += showBinds(m.keybinds)
	return s
}
