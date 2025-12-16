package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type DevboxConfig struct {
	Schema   string      `json:"$schema"`
	Packages []string    `json:"packages"`
	Shell    ShellConfig `json:"shell"`
}

type ShellConfig struct {
	InitHook []string            `json:"init_hook"`
	Scripts  map[string][]string `json:"scripts"`
}

type ScriptsConfig struct {
	Test []string `json:"test"`
}

type Component interface {
	Compatible() bool
	Name() string
}

type Project struct {
	Path   string
	Config DevboxConfig
}

func NewProject(p ...string) Project {
	var path string
	var err error
	if len(p) == 0 {
		path, err = os.Getwd()
		if err != nil {
			// fallback if getwd fails, though unlikely
			path = "."
		}
	} else {
		path = p[0]
	}
	devboxPath := filepath.Join(path, "devbox.json")
	file, err := os.Open(devboxPath)
	if err != nil {
		return Project{Path: path}
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var config DevboxConfig
	json.Unmarshal(byteValue, &config)

	return Project{Path: path, Config: config}
}

func (p *Project) Compatible() bool {
	devboxPath := filepath.Join(p.Path, "devbox.json")
	if _, err := os.Stat(devboxPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Name returns the name of the project, which is the last part of the path.
func (p *Project) Name() string {
	return filepath.Base(p.Path)
}

// Commands returns the commands defined in shell.scripts section of the devbox.json file.
func (p *Project) Commands() map[string][][]string {
	convertedScripts := make(map[string][][]string)
	for key, commands := range p.Config.Shell.Scripts {
		var parsedCommands [][]string
		for _, command := range commands {
			parsedCommands = append(parsedCommands, strings.Split(command, " "))
		}
		convertedScripts[key] = parsedCommands
	}
	return convertedScripts
}

func (p *Project) CommandNames() []string {
	var names []string
	for k := range p.Commands() {
		names = append(names, k)
	}
	return names
}

// Info returns the project name and the commands defined in the devbox.json file.
func (p *Project) Info() []string {
	output := []string{fmt.Sprintf("Project name: %s", p.Name())}
	output = append(output, "Commands available through 'camp project [command]':")
	for _, c := range p.CommandNames() {
		output = append(output, " - "+c)
	}
	return output
}
