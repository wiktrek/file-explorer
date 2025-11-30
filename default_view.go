package main

import "fmt"

func defaultView(m model) string {
	s := "Select File: \n\n"
	if !m.config.hidePath {
		s += m.currentDir + "\n"
	}

	for i, file := range m.files {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, file.fileType, file.path)
	}
	if m.searching {
		if m.secondCursor == len(m.search) {
			s += "/" + m.search + "|"
		} else {
			s += "/"
			for j := range m.search {
				if j == m.secondCursor {
					s += fmt.Sprintf("|%c", m.search[j])
				} else {
					s += fmt.Sprintf("%c", m.search[j])
				}
			}
		}
	} else if m.search != "" {
		s += "/" + m.search
	}
	s += showBinds(m.config.keybinds)
	return s
}
