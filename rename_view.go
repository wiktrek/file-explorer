package main

import "fmt"

func renameView(m model) string {
	s := "Rename Your files: "
	if !m.config.hidePath {
		s += m.currentDir + "\n"
	}

	for i, file := range m.files {
		if m.cursor == i {
			//  handle renaming
			s += "> "

			if m.secondCursor == len(m.tempString) {
				s += m.tempString
				s += "|"
			} else {
				for j := range m.tempString {
					if j == m.secondCursor {
						s += fmt.Sprintf("|%c", m.tempString[j])
					} else {
						s += fmt.Sprintf("%c", m.tempString[j])
					}
				}
			}
			s += " <\n"
		} else {
			s += fmt.Sprintf("%s %s\n", file.fileType, file.path)
		}
	}
	s += showBinds(m.config.keybinds)
	return s
}
