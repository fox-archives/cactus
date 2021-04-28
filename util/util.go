package util

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	g "github.com/AllenDang/giu"
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
