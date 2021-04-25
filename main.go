package main

import (
	"image/color"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	cli "github.com/urfave/cli/v2"
)

type GlobalCmdResult struct {
	err         error
	output      string
	cfgEntryKey string
	cfgEntry    CfgEntry
}

var globalCmdResult = &GlobalCmdResult{
	err:         nil,
	output:      "",
	cfgEntryKey: "",
	cfgEntry:    CfgEntry{},
}

func loop(cfgCactus CfgCactus, binds CfgBind) {
	if g.IsKeyDown(g.KeyEscape) {
		os.Exit(0)
	}

	// key is the highest parent properties of the config,
	// who's value is cfgEntry
	for key, bindEntry := range binds {
		if g.IsKeyDown(keyMap[key]) {
			output, didRun, err := runCmdOnce(key, bindEntry)
			if didRun {
				if err != nil {
					globalCmdResult = &GlobalCmdResult{
						err:         err,
						output:      output,
						cfgEntryKey: key,
						cfgEntry:    bindEntry,
					}
				} else {
					os.Exit(0)
				}
			}
		}
	}

	var widgets []g.Widget
	// if there is some error, prepend the output
	if globalCmdResult.err != nil {
		widgets = append(widgets, g.Label("Error"))
		widgets = append(widgets, g.Label(globalCmdResult.err.Error()))

		widgets = append(widgets, g.Label("Output"))
		widgets = append(widgets, g.Label(globalCmdResult.output))

		table := g.Table("Command Table").Rows(g.Row(
			g.Label(globalCmdResult.cfgEntryKey),
			g.Label(globalCmdResult.cfgEntry.Cmd),
		).BgColor(&(color.RGBA{200, 100, 100, 255})))
		widgets = append(widgets, table)
	} else {
		table := g.Table("Command Table").FastMode(true).Rows(buildGuiTableRows(binds)...)
		widgets = append(widgets, table)
	}

	g.SingleWindow("Runner").Layout(
		widgets...,
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
				Name:    "binds",
				Aliases: []string{"b"},
				Value:   getCfgFile("binds.ini"),
				Usage:   "Location of bindings",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   getCfgFile("cactus.toml"),
				Usage:   "Location of configuration file",
			},
		},
		Action: func(c *cli.Context) error {
			binds := getCfgBinds(c.String("binds"))
			cfg := getCfgCactus(c.String("config"))

			imgui.CreateContext(nil)

			if cfg.FontFile != "" {
				imgui.CurrentIO().Fonts().AddFontFromFileTTF(cfg.FontFile, float32(cfg.FontSize))
			}

			wnd := g.NewMasterWindow("Cactus", 750, 450, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating, nil)

			wnd.Run(func() {
				loop(cfg, binds)
			})

			return nil
		},
	}

	err := app.Run(os.Args)
	handle(err)
}
