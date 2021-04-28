package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/alessio/shellescape"
	"github.com/eankeen/cactus/cfg"
)

func Handle(err error) {
	if err != nil {
		_, isDebug := os.LookupEnv("DEBUG")
		if isDebug {
			panic(err)
		} else {
			log.Fatalln(err)
		}
	}
}

type SystemdRunOutput struct {
	RunningAsUnit             string
	FinishedWithResult        string
	MainProcessTerminatedWith string
	ServiceRuntime            string
	CPUTimeConsumed           string
}

func ParseSystemdRunOutput(output string) [][]string {
	// An array of key value pairs
	keyValueArr := [][]string{}

	for _, line := range strings.Split(output, "\n") {
		arr := strings.Split(line, ":")
		if len(arr) < 2 {
			continue
		}

		key := arr[0]
		value := arr[1]
		value = strings.TrimSpace(value)

		keyValueArr = append(keyValueArr, []string{key, value})
	}

	return keyValueArr
}

func CopyToClipboard(data string) {
	data = strings.TrimSpace(data)

	var cmd *exec.Cmd

	_, err := os.Stat("/usr/bin/dash")
	// TODO: check if xclip is not installed (show error at top of screen, along with other "Internal Errors?")
	if errors.Is(err, os.ErrNotExist) {
		cmd = exec.Command("/usr/bin/sh", "-c", fmt.Sprintf("echo %s | xclip -r -selection clipboard", shellescape.Quote(data)))
	} else if err != nil {
		Handle(err)
	} else {
		cmd = exec.Command("/usr/bin/dash", "-c", fmt.Sprintf("echo %s | xclip -r -selection clipboard", shellescape.Quote(data)))
	}

	err = cmd.Start()
	Handle(err)
}

// Build array of rows for the main display
func BuildGuiTableRows(keybinds cfg.Keybinds) []*g.RowWidget {
	var sortedKeys = make([]string, 0, len(keybinds))
	for key := range keybinds {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// Create the UI rows
	var rowWidgets = make([]*g.RowWidget, 0, len(keybinds))
	for _, key := range sortedKeys {
		value := keybinds[key]

		rowWidgets = append(rowWidgets, g.Row(
			g.Label(key),
			g.Label(value.Cmd+" "+strings.Join(value.Args, " ")),
		))
	}

	return rowWidgets
}

// Get a particular file path from the configuration directory
func GetCfgFile(file string) string {
	var cfgDir string

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		cfgDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(cfgDir, "cactus", file)
}
