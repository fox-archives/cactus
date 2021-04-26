package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	g "github.com/AllenDang/giu"
	"github.com/eankeen/cactus/cfg"
	"github.com/google/uuid"
)

var hasRan = false

func handle(err error) {
	if err != nil {
		_, isDebug := os.LookupEnv("DEBUG")
		if isDebug {
			panic(err)
		} else {
			fmt.Println(err)
		}
	}
}

func runCmd(key string, keyBind cfg.KeybindEntry) (string, bool, error) {
	var cmd *exec.Cmd

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", false, fmt.Errorf("Error: Could not generate random number\n%w", err)
	}

	args := []string{
		"--no-ask-password",
		"--unit", "cactus-" + uuid.String(),
		"--description", fmt.Sprintf("Cactus Start for command: '%s'", keyBind.Cmd),
		"--send-sighup",
		"--working-directory",
		os.Getenv("HOME"), "--user",
	}

	switch keyBind.Run {
	case "dash":
		args = append(args, "/usr/bin/dash", "-c", keyBind.Cmd)
	case "bash":
		args = append(args, "/usr/bin/bash", "-c", keyBind.Cmd)
	default:
		args = append(args, keyBind.Cmd)
	}

	cmd = exec.Command("/usr/bin/systemd-run", args...)

	output, err := cmd.CombinedOutput()
	return string(output), true, err
}

func runCmdOnce(key string, keyBind cfg.KeybindEntry) (string, bool, error) {
	if hasRan == true {
		return "", false, nil
	}

	hasRan = true
	return runCmd(key, keyBind)
}

func buildGuiTableRows(keyBinds cfg.Keybinds) []*g.RowWidget {
	// created sort keys
	// without this, the menu will ranomdly change
	// ordering, as go does
	var sortedKeys = make([]string, 0, len(keyBinds))
	for key := range keyBinds {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// create the ui rows
	var rowWidgets = make([]*g.RowWidget, 0, len(keyBinds))
	for _, key := range sortedKeys {
		value := keyBinds[key]

		rowWidgets = append(rowWidgets, g.Row(
			g.Label(key),
			g.Label(value.Cmd),
		))
	}

	return rowWidgets
}

// Get a particular file path from the configuration directory
func getCfgFile(file string) string {
	var cfgDir string

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		cfgDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(cfgDir, "cactus", file)
}
