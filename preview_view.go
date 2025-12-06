package main

import "fmt"

func previewView(m model) string {
	var s string
	path := m.currentDir + m.files[m.cursor].path
	o, err := IsDirectory(path)
	if err != nil {
		fmt.Printf("%v", err)
		return ""
	}
	if o {
		s += "dir"
		files := loadFiles(path)
		for _, file := range files {
			s += fmt.Sprintf("%s %s\n", file.fileType, file.path)
		}
	} else {
		s += "file"
		s += readFile(path)
	}
	// s += readFile(m.currentDir + m.files[m.cursor].path)
	return s
}
