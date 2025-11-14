package main

type Keybind = struct {
	keybind     string
	description string
}
type File = struct {
	fileType string
	path     string
}

var keybinds = []Keybind{
	{
		keybind:     "ctrl+q",
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
	{
		keybind:     "ctrl+n",
		description: "new file",
	},
}

type ViewState string

const (
	Default       ViewState = "default"
	ConfirmDelete ViewState = "confirm"
	Rename        ViewState = "rename"
	Move          ViewState = "move"
	Copy          ViewState = "copy"
	New           ViewState = "New"
)

type model struct {
	files        []File
	cursor       int
	selected     map[int]struct{}
	currentDir   string
	secondCursor int
	temp_string  string
	config       Config
	viewState    ViewState
}

var config = getConfig("config.json")
