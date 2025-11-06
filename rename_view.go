package main

import "fmt"

func renameView(m model) string {
	s := "Rename Your files"
	s += m.currentDir + "\n"

	for i, file := range m.files {
		if m.cursor == i {
			//  handle renaming
			s += "> "

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
			s += " <\n"
		} else {
			s += fmt.Sprintf("%s\n", file)
		}
	}
	s += showBinds(m.keybinds)
	return s
}
