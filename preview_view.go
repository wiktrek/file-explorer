package main

import (
	"fmt"
	"math"
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

	maxWidth := int(math.Floor(float64(m.width)/float64(m.viewCount))) - 2
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
			if len(line) > maxWidth {
				line = line[:maxWidth-3] + "..."
			}
			s += line + "\n"
		}
	}
	// s += readFile(m.currentDir + m.files[m.cursor].path)
	return s
}
