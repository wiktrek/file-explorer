package main

import (
	"fmt"
)

func moveView(m model) string {
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
	// split := strings.Split(m.tempString, "/")
	// file_name := split[len(split)-1]
	// s += " _" + file_name + "\n"
	s += "Moving " + m.tempString
	s += "\nPress p to switch do default"
	return s
}
