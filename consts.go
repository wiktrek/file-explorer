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
		keybind:     "ctrl+x/m",
		description: "move file",
	},
	{
		keybind:     "ctrl+c",
		description: "copy file",
	},
	{
		keybind:     "ctrl+v",
		description: "paste file",
	},
}

type ViewState string

const (
	Default       ViewState = "default"
	ConfirmDelete ViewState = "confirm"
	Rename        ViewState = "rename"
	Move          ViewState = "move"
	Copy          ViewState = "copy"
)

type model struct {
	files        []string
	cursor       int
	selected     map[int]struct{}
	currentDir   string
	secondCursor int
	temp_string  string
	keybinds     bool
	viewState    ViewState
}

var config = getConfig("config.json")
