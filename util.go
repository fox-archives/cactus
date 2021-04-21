package main

import (
	"os"
	"sort"
	"syscall"

	g "github.com/AllenDang/giu"
)

var isRunning = false

func runCmdOnce(key string, cmd string, runWith string) {
	if isRunning == true {
		return
	}

	isRunning = true

	// stdin, err := os.Open("/dev/null")
	// handle(err)
	// stdout, err := os.OpenFile("/dev/null", os.O_RDWR, os.ModeCharDevice)
	// handle(err)
	// stderr, err := os.OpenFile("/dev/null", os.O_RDWR, os.ModeCharDevice)
	// handle(err)

	var sysproc = &syscall.SysProcAttr{Noctty: true}
	var attr = os.ProcAttr{
		Dir: os.Getenv("HOME"),
		Env: os.Environ(),
		Sys: sysproc,
	}

	var name string
	var argv []string

	if runWith == "" {
		name = "/usr/bin/systemd-run"
		argv = []string{"systemd-run", "--user", "bash", "-c", cmd}
	} else {
		// TODO: pass into exec manually
		name = "/bin/bash"
		argv = []string{name, "-c", cmd}
	}

	process, err := os.StartProcess("systemd-run", argv, &attr)
	// process, err := os.StartProcess(name, argv, &attr)
	handle(err)

	err = process.Release()
	handle(err)
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
