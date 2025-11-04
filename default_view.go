package main

import "fmt"

func showBinds() string {
	s := ""

	for k := range keybinds {
		s += "\nPress " + keybinds[k].keybind + " to " + keybinds[k].description
	}
	return s
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
	s += showBinds()
	return s
}
