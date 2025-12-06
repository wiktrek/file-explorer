package main

import "fmt"

func newView(m model) string {
	s := "Create new file: \n\n"
	if !m.config.hidePath {
		s += m.currentDir + "\n"
	}

	for _, file := range m.files {
		s += fmt.Sprintf(" %s %s\n", file.fileType, file.path)
	}
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
	s += "\n"
	return s
}
