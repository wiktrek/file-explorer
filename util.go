package main

import "fmt"

func remove_at_index(str string, i int) string {
	if i == 0 {
		return str
	}
	if i == len(str) {
		return str[:i-1]
	}
	return str[:i] + str[(i+1):]
}
func add_to_string(str string, c byte, i int) string {
	if i == len(str) {
		return fmt.Sprintf("%s%c", str, c)
	}
	return fmt.Sprintf("%s%c%s", str[:i], c, str[i:])
}
