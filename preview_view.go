package main

import (
	"fmt"
	"strings"
)

func previewView(m model) string {
	var s string
	path := m.currentDir + m.files[m.cursor].path
	o, err := IsDirectory(path)
	if err != nil {
		fmt.Printf("%v", err)
		return ""
	}

	s += m.files[m.cursor].path + ":\n\n"
	if o {
		files := loadFiles(path + "/")
		for i, file := range files {
			if i >= m.height-10 {
				s += "..."
				break
			}
			s += fmt.Sprintf("%s %s\n", file.fileType, file.path)
		}
	} else {

		lines := strings.Split(readFile(path), "\n")
		for i, line := range lines {
			if i >= m.height-20 {
				s += "..."
				break
			}
			s += line + "\n"
		}
	}
	// s += readFile(m.currentDir + m.files[m.cursor].path)
	return s
}
