package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type JsonConfig struct {
	DefaultPath string
	HidePath    bool
	Keybinds    bool
}
type Config struct {
	defaultPath string
	hidePath    bool
	keybinds    bool
}

func getConfig(configDir string) Config {
	var config JsonConfig
	b, err := os.ReadFile(configDir)
	if err != nil {
		fmt.Println(err)
	}
	json.NewDecoder(bytes.NewBuffer(b)).Decode(&config)
	if !pathExits(config.DefaultPath) {
		config.DefaultPath = "/home/"
	}
	// I hate the uppercase naming ( I will probably change this later)
	return Config{
		defaultPath: config.DefaultPath,
		hidePath:    config.HidePath,
		keybinds:    config.Keybinds,
	}
}
func loadFiles(dir string) []File {
	pattern := "*"
	matches, err := filepath.Glob(dir + pattern)
	if err != nil {
		log.Fatal(err)
	}
	var files []File
	for _, match := range matches {
		split := strings.Split(match, "/")
		file_name := split[len(split)-1]
		files = append(files, File{
			path:     file_name,
			fileType: getIcon(dir + file_name),
		})
	}
	return files
}
func goUp(dir string) string {

	if dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	var split = strings.Split(dir, "/")
	result := strings.Join(split[:len(split)-1], "/")

	if len(result) <= 1 {
		return "/"
	}
	return result + "/"
}

func reloadDir(m model, path string) model {
	m.cursor = 0
	if path != "" {
		m.currentDir = path
	}
	m.files = loadFiles(m.currentDir)
	m.View()
	return m
}
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
func openFile(f string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, f)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
func pathExits(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func deletePath(path string) {
	if path == "/" {
		return
	}
	os.Remove(path)
}
func moveFile(path string, path2 string) error {
	err := os.Rename(path, path2)
	return err
}
func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
func createFile(path string) {
	os.Create(path)
}
