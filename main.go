package main

import (
	"fmt"
	"os"
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
		argv = []string{"-c", cmd}
	} else {
		// TODO: pass into exec manually
		name = "/bin/bash"
		argv = []string{"-c", cmd}
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
		g.Label("Thing"),
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
