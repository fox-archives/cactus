package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	g "github.com/AllenDang/giu"
	"github.com/google/uuid"
)

var isRunning = false

func runCmd(key string, cfgEntry CfgEntry) (string, bool, error) {
	var cmd *exec.Cmd

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", false, err
	}

	args := []string{
		"--no-ask-password",
		"--unit", "cactus-" + uuid.String(),
		"--description", fmt.Sprintf("Cactus Start for command: '%s'", cfgEntry.Cmd),
		"--send-sighup",
		"--working-directory",
		os.Getenv("HOME"), "--user",
	}

	switch cfgEntry.Run {
	case "dash":
		args = append(args, "/usr/bin/dash", "-c", cfgEntry.Cmd)
	case "bash":
		args = append(args, "/usr/bin/bash", "-c", cfgEntry.Cmd)
	default:
		args = append(args, cfgEntry.Cmd)
	}

	cmd = exec.Command("/usr/bin/systemd-run", args...)

	output, err := cmd.CombinedOutput()
	return string(output), true, err
}

func runCmdOnce(key string, cfgEntry CfgEntry) (string, bool, error) {
	if isRunning == true {
		return "", false, nil
	}

	isRunning = true
	return runCmd(key, cfgEntry)
}

func buildGuiTableRows(cfgToml CfgToml) []*g.RowWidget {
	// created sort keys
	// without this, the menu will ranomdly change
	// ordering, as go does
	var sortedKeys = make([]string, 0, len(cfgToml))
	for key := range cfgToml {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// create the ui rows
	var rowWidgets = make([]*g.RowWidget, 0, len(cfgToml))
	for _, key := range sortedKeys {
		value := cfgToml[key]

		rowWidgets = append(rowWidgets, g.Row(
			g.Label(key),
			g.Label(value.Cmd),
		))
	}

	return rowWidgets
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
