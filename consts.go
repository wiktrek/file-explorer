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
	},
	{
		keybind:     "d",
		description: "delete file",
	},
	{
		keybind:     "f2",
		description: "rename file",
	},
	{
		keybind:     "m",
		description: "move file",
	},
}

type ViewState string

const (
	Default       ViewState = "default"
	ConfirmDelete ViewState = "confirm"
	Rename        ViewState = "rename"
	Move          ViewState = "move"
)

type model struct {
	files        []string
	cursor       int
	selected     map[int]struct{}
	currentDir   string
	secondCursor int
	temp_string  string
	viewState    ViewState
}

var defaultDir = "/home/wiktor/projects/golang/file-explorer/"
