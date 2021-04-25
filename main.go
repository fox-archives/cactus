package main

import (
	"os"

	g "github.com/AllenDang/giu"
	cli "github.com/urfave/cli/v2"
)

func loop(cfg CfgToml) {
	if g.IsKeyDown(g.KeyEscape) {
		os.Exit(0)
	}

	// key is the highest parent properties of the config,
	// who's value is cfgEntry
	for key, cfgEntry := range cfg {
		if g.IsKeyDown(keyMap[key]) {
			runCmdOnce(key, cfgEntry)
			os.Exit(1)
		}
	}

	g.SingleWindow("Runner").Layout(
		g.Table("Command Table").FastMode(true).Rows(buildGuiTableRows(cfg)...),
	)
}

func main() {
	app := &cli.App{
		Name:    "cactus",
		Usage:   "Small hotkey application",
		Version: "0.3.0",
		Authors: []*cli.Author{
			{
				Name:  "Edwin Kofler",
				Email: "edwin@kofler.dev",
			},
		},
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   getCfgFile(),
				Usage:   "Location of configuration file",
			},
		},
		Action: func(c *cli.Context) error {
			wnd := g.NewMasterWindow("Cactus", 750, 450, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating, nil)

			wnd.Run(func() {
				loop(getCfg(c.String("config")))
			})

			return nil
		},
	}

	err := app.Run(os.Args)
	handle(err)
}
