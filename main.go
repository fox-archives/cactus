package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	g "github.com/AllenDang/giu"
	cli "github.com/urfave/cli/v2"
)

var isRunning = false

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func buildCmdMenu(cfgToml CfgToml) []*g.RowWidget {
	// created sort keys
	// without this, the menu will ranomdly change
	// ordering, as go does
	var sortedKeys = make([]string, 0, len(cfgToml))
	for key := range cfgToml {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// create the ui
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

func runCmdOnce(key string, cmd string, runWith string) {
	if isRunning == true {
		return
	}

	isRunning = true

	var sysproc = &syscall.SysProcAttr{Noctty: true}
	var attr = os.ProcAttr{
		Dir: os.Getenv("HOME"),
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			nil,
			nil,
		},
		Sys: sysproc,
	}

	var name string
	var argv []string

	if runWith == "" {
		name = "/bin/bash"
		argv = []string{name, "-c", cmd}
	} else {
		// TODO: pass into exec manually
		name = "/bin/bash"
		argv = []string{name, "-c", cmd}
	}

	process, err := os.StartProcess(name, argv, &attr)
	handle(err)

	err = process.Release()
	handle(err)
}

func loop(cfg CfgToml) {
	for key, value := range cfg {
		if g.IsKeyDown(keyMap[key]) {
			runCmdOnce(key, value.Cmd, value.Run)
			os.Exit(1)
		}
	}

	g.SingleWindow("Runner").Layout(
		g.Table("Command Table").FastMode(true).Rows(buildCmdMenu(cfg)...),
	)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
		Action: func(c *cli.Context) error {
			name := "Nefertiti"
			if c.NArg() > 0 {
				name = c.Args().Get(0)
			}
			if c.String("lang") == "spanish" {
				fmt.Println("Hola", name)
			} else {
				fmt.Println("Hello", name)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	handle(err)

	// TODO: pass --config in cli
	cfg := getCfg()

	wnd := g.NewMasterWindow("Cactus", 750, 450, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating, nil)

	wnd.Run(func() {
		loop(cfg)
	})
}
