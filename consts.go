package main

type Keybind = struct {
	keybind     string
	description string
}

var keybinds = []Keybind{
	{
		keybind:     "q",
		description: "quit",
	},
	{
		keybind:     "o",
		description: "open file",
	}, {
		keybind:     "d",
		description: "delete file",
	},
}

type ViewState string

const (
	Default       ViewState = "default"
	ConfirmDelete ViewState = "confirm"
)

type model struct {
	files         []string
	cursor        int
	selected      map[int]struct{}
	currentDir    string
	confirmCursor int
	viewState     ViewState
}

var defaultDir = "/home/wiktor/projects/golang/file-explorer/"
