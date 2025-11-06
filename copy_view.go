package main

import "fmt"

func copyView(m model) string {
	s := "Select File: \n\n"
	s += m.currentDir + "\n"

	for i, file := range m.files {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, file)
	}
	s += fmt.Sprintf("Copying file: %s\n", m.temp_string)
	s += showBinds(m.keybinds)
	return s
}
