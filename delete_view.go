package main

func confirmDeleteView(m model) string {
	s := "Are you really sure you want to delete "
	s += m.currentDir + m.files[m.cursor].path + "\n"
	if m.secondCursor == 0 {
		s += ">"
	}
	s += "YES"
	s += "    "
	if m.secondCursor == 1 {
		s += ">"
	}
	s += "NO"
	return s
}
