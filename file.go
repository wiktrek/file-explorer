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
		home, _ := os.UserHomeDir()
		config.DefaultPath = home
	}
	// I hate the uppercase naming ( I will probably change this later)
	return Config{
		defaultPath: config.DefaultPath,
		hidePath:    config.HidePath,
		keybinds:    config.Keybinds,
	}
}
func loadFiles(dir string) []File {
	if !strings.HasSuffix(dir, string(filepath.Separator)) {
		dir += string(filepath.Separator)
	}
	pattern := "*"
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		log.Fatal(err)
	}
	var files []File
	for _, match := range matches {
		file_name := filepath.Base(match)
		files = append(files, File{
			path:     file_name,
			fileType: getIcon(filepath.Join(dir, file_name)),
		})
	}
	return files
}
func goUp(dir string) string {
	if !strings.HasSuffix(dir, string(filepath.Separator)) {
		dir = dir[:len(dir)-1]
	}
	parent := filepath.Dir(dir)
	if parent == dir {
		return dir
	}
	if !strings.HasSuffix(parent, string(filepath.Separator)) {
		parent += string(filepath.Separator)
	}
	return parent
}

func reloadDir(m model, path string) model {
	m.cursor = 0
	if path != "" {
		m.currentDir = path
	}
	m.files = search_filter(loadFiles(m.currentDir), m.search)
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
