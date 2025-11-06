package main

import "fmt"

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
	s += showBinds(m.keybinds)
	return s
}
