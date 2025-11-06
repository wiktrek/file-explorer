package main

import "fmt"

func remove_at_index(str string, i int) string {
	if i == 0 {
		return str
	}
	if i == len(str) {
		return str[:i-1]
	}
	return str[:i-1] + str[(i):]
}
func add_to_string(str string, c byte, i int) string {
	if i == len(str) {
		return fmt.Sprintf("%s%c", str, c)
	}
	return fmt.Sprintf("%s%c%s", str[:i], c, str[i:])
}
func showBinds(on bool) string {
	s := ""
	if on {
		for k := range keybinds {
			s += "\nPress " + keybinds[k].keybind + " to " + keybinds[k].description
		}
	}
	return s
}
func getIcon(file string) string {
	r, err := IsDirectory(file)
	if err != nil {
		fmt.Println(err)
	}
	if r {
		return "F"
	} else {
		return "P"
	}
}
