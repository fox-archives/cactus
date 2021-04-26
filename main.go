package main

import (
	"os"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/eankeen/cactus/cfg"
	"github.com/fsnotify/fsnotify"
	cli "github.com/urfave/cli/v2"
)

type GlobalCmdResult struct {
	err        error
	output     string
	keyBindKey string
	keyBind    cfg.KeybindEntry
}

var globalCmdResult = &GlobalCmdResult{
	err:        nil,
	output:     "",
	keyBindKey: "",
	keyBind:    cfg.KeybindEntry{},
}

func runCmdOnceWrapper(key string, bindEntry cfg.KeybindEntry) {
	output, didRun, err := runCmdOnce(key, bindEntry)
	if didRun {
		if err != nil {
			globalCmdResult = &GlobalCmdResult{
				err:        err,
				output:     output,
				keyBindKey: key,
				keyBind:    bindEntry,
			}
		} else {
			os.Exit(0)
		}
	}
}

func loop(cfg *cfg.Cfg, binds *cfg.Keybinds) {
	if g.IsKeyDown(g.KeyEscape) {
		os.Exit(0)
	}

	// key is the highest parent properties of the config,
	// who's value is cfgEntry
	for key, bindEntry := range *binds {
		var mod string

		if strings.Contains(key, "-") {
			strs := strings.Split(key, "-")
			mod = strs[0]
			key = strs[1]

			if mod == "Shift" && (g.IsKeyDown(g.KeyLeftShift) || g.IsKeyDown(g.KeyRightShift)) && g.IsKeyDown(keyMap[key]) {
				runCmdOnceWrapper(key, bindEntry)
			} else if mod == "Control" && (g.IsKeyDown(g.KeyLeftControl) || g.IsKeyDown(g.KeyRightControl)) && g.IsKeyDown(keyMap[key]) {
				runCmdOnceWrapper(key, bindEntry)
			} else if mod == "Alt" && (g.IsKeyDown(g.KeyLeftAlt) || g.IsKeyDown(g.KeyRightAlt)) && g.IsKeyDown(keyMap[key]) {
				runCmdOnceWrapper(key, bindEntry)
			}
		} else {
			if mod == "" && g.IsKeyDown(keyMap[key]) {
				runCmdOnceWrapper(key, bindEntry)
			}
		}

	}

	var widgets []g.Widget
	// if there is some error, show it to the top
	if globalCmdResult.err != nil {
		widgets = append(widgets, g.Label("ERROR"))
		widgets = append(widgets, g.Label(globalCmdResult.err.Error()))
		widgets = append(widgets, g.Label(""))

		widgets = append(widgets, g.Label("OUTPUT"))
		widgets = append(widgets, g.Label(globalCmdResult.output))
		widgets = append(widgets, g.Label(""))

		widgets = append(widgets, g.Label("KEY"))
		widgets = append(widgets, g.Label("Mod: "+globalCmdResult.keyBind.Mod))
		widgets = append(widgets, g.Label("Run: "+globalCmdResult.keyBind.Run))
		widgets = append(widgets, g.Label("Cmd: "+globalCmdResult.keyBind.Cmd))
	} else {
		table := g.Table("Command Table").FastMode(true).Rows(buildGuiTableRows(*binds)...).Flags(
			imgui.TableFlags_Resizable | imgui.TableFlags_RowBg | imgui.TableFlags_Borders | imgui.TableFlags_SizingFixedFit | imgui.TableFlags_ScrollX | imgui.TableFlags_ScrollY,
		)
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
				Value:   getCfgFile("binds.toml"),
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
			// configuration
			keybindsMnger := cfg.NewKeybindsMnger(c.String("binds"))
			err := keybindsMnger.Reload()
			handle(err)

			cfgMnger := cfg.NewCfgMnger(c.String("config"))
			err = cfgMnger.Reload()
			handle(err)

			cfg := cfgMnger.Get()
			keybindings := keybindsMnger.Get()

			// watcher
			watcher, err := fsnotify.NewWatcher()
			handle(err)
			defer watcher.Close()

			done := make(chan bool)
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}

						if event.Op&fsnotify.Write == fsnotify.Write {
							err = keybindsMnger.Reload()
							handle(err)

							err = cfgMnger.Reload()
							handle(err)

							// TODO ?
							// imgui.NewFrame()
							// imgui.Render()
							// drawData := imgui.RenderedDrawData()
							// fmt.Println(drawData.Valid())
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}

						handle(err)
					}
				}
			}()

			err = watcher.Add(c.String("binds"))
			handle(err)
			err = watcher.Add(c.String("config"))
			handle(err)

			// imgui
			ctx := imgui.CreateContext(nil)
			err = ctx.SetCurrent()
			handle(err)

			if cfg.FontFile != "" {
				imgui.CurrentIO().Fonts().AddFontFromFileTTF(cfg.FontFile, float32(cfg.FontSize))
			}
			// imgui.CurrentIO().Fonts().Build()

			wnd := g.NewMasterWindow("Cactus", 800, 450, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating, nil)

			wnd.Run(func() {
				loop(cfg, keybindings)
			})

			<-done
			return nil
		},
	}

	err := app.Run(os.Args)
	handle(err)
}
