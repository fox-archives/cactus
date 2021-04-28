package main

import (
	"fmt"
	"os"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/eankeen/cactus/cfg"
	cmd "github.com/eankeen/cactus/cmd"
	"github.com/eankeen/cactus/keymap"
	"github.com/eankeen/cactus/util"
	"github.com/fsnotify/fsnotify"
	cli "github.com/urfave/cli/v2"
)

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
				Value:   util.GetCfgFile("binds.toml"),
				Usage:   "Location of bindings",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   util.GetCfgFile("cactus.toml"),
				Usage:   "Location of configuration file",
			},
		},
		Action: func(c *cli.Context) error {
			/* -------------------- Configuration ------------------- */
			keybindsMnger := cfg.NewKeybindsMnger(c.String("binds"))
			err := keybindsMnger.Reload()
			util.Handle(err)

			cfgMnger := cfg.NewCfgMnger(c.String("config"))
			err = cfgMnger.Reload()
			util.Handle(err)

			cfg := cfgMnger.Get()
			keybinds := keybindsMnger.Get()

			/* ----------------------- Watcher ---------------------- */
			watcher, err := fsnotify.NewWatcher()
			util.Handle(err)
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
							util.Handle(err)

							err = cfgMnger.Reload()
							util.Handle(err)

							// TODO: Manually rerender?
							// imgui.NewFrame()
							// imgui.Render()
							// drawData := imgui.RenderedDrawData()
							// fmt.Println(drawData.Valid())
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}

						util.Handle(err)
					}
				}
			}()

			err = watcher.Add(c.String("binds"))
			util.Handle(err)
			err = watcher.Add(c.String("config"))
			util.Handle(err)

			/* ------------------------ Imgui ----------------------- */
			ctx := imgui.CreateContext(nil)
			err = ctx.SetCurrent()
			util.Handle(err)

			if cfg.FontFile != "" {
				imgui.CurrentIO().Fonts().AddFontFromFileTTF(cfg.FontFile, float32(cfg.FontSize))
			}

			wnd := g.NewMasterWindow("Cactus", 800, 450, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating, nil)

			// The key the user wants to run. By default, it's blank
			var myCmd = cmd.New()
			wnd.Run(func() {
				loop(cfg, keybinds, myCmd)
			})

			<-done
			return nil
		},
	}

	err := app.Run(os.Args)
	util.Handle(err)
}

func loop(cfg *cfg.Cfg, keybinds *cfg.Keybinds, myCmd *cmd.Cmd) {
	// Quit on Escape
	if g.IsKeyDown(g.KeyEscape) {
		os.Exit(0)
	}

	// As we iterate over *keybinds, we cannot execute immediately
	// because there could be a more specific keybind later in the
	// iteration. Here, we save information  about matches. It has a
	// format like 'G' or 'Ctrl-G'
	matched := ""

	// 'key' is each key in the config file,
	// The properties of each key is enumerated in the
	// members of keybindEntry
	for key := range *keybinds {
		// If there is no hypthen in a config key, it contains a modifier
		if !strings.Contains(key, "-") {
			if g.IsKeyDown(keymap.Keymap[key]) {
				// Only overridie matched if doesn't contain a hypthen (which implies
				// it is a more specific modifier)
				if !strings.Contains(key, "-") {
					matched = key
				}
			}
		} else {
			// If there is a hypthen, take out the modifier key
			arr := strings.Split(key, "-")
			modifier := arr[0]
			actualKey := arr[1]

			if modifier == "Shift" && (g.IsKeyDown(g.KeyLeftShift) || g.IsKeyDown(g.KeyRightShift)) && g.IsKeyDown(keymap.Keymap[actualKey]) {
				matched = key
			} else if (modifier == "Ctrl" || modifier == "Control") && (g.IsKeyDown(g.KeyLeftControl) || g.IsKeyDown(g.KeyRightControl)) && g.IsKeyDown(keymap.Keymap[actualKey]) {
				matched = key
			} else if (modifier == "Alt") && (g.IsKeyDown(g.KeyLeftAlt) || g.IsKeyDown(g.KeyRightAlt)) && g.IsKeyDown(keymap.Keymap[actualKey]) {
				matched = key
			}
		}
	}

	// If we have a match, attempt to execute the keybinding so long as
	// we haven't already done so
	if !myCmd.HasRan && matched != "" {
		fmt.Println("before", !myCmd.Keybind.AlwaysShowInfo)

		myCmd.HasRan = true

		if strings.Contains(matched, "-") {
			arr := strings.Split(matched, "-")
			modifier := arr[0]
			key := arr[1]

			myCmd.KeybindMod = modifier
			myCmd.KeybindKey = key
			myCmd.Keybind = (*keybinds)[matched]
			myCmd.Result = myCmd.RunCmd()
		} else {
			myCmd.KeybindMod = ""
			myCmd.KeybindKey = matched
			myCmd.Keybind = (*keybinds)[matched]
			myCmd.Result = myCmd.RunCmd()
		}

		// Exit if there is a success and we don't want to show info on success
		fmt.Println(myCmd.Result.Err == nil, !myCmd.Keybind.AlwaysShowInfo)
		if myCmd.Result.Err == nil && !myCmd.Keybind.AlwaysShowInfo {
			os.Exit(0)
		}
	}

	// SHOW THE GUI
	var widgets []g.Widget
	if !myCmd.HasRan {
		// If no commands were run, show the table
		table := g.Table("Command Table").FastMode(true).Rows(util.BuildGuiTableRows(*keybinds)...).Flags(
			imgui.TableFlags_Resizable | imgui.TableFlags_RowBg | imgui.TableFlags_Borders | imgui.TableFlags_SizingFixedFit | imgui.TableFlags_ScrollX | imgui.TableFlags_ScrollY | imgui.TableFlags_ScrollY,
		)
		widgets = append(widgets, table)
	} else {
		// If we are this far and the commands are

		/* ----------------------- RESULT ----------------------- */
		if myCmd.Result.Err != nil {
			widgets = append(widgets, g.Line(
				g.ArrowButton("Arrow", g.DirectionRight),
				g.Label("RESULT"),
			))

			widgets = append(widgets, g.Label("Error: "+myCmd.Result.Err.Error()))
			widgets = append(widgets, g.Label("ExecName: "+myCmd.Result.ExecName))
			widgets = append(widgets, g.Label("ExecArgs: ["))
			for _, arg := range myCmd.Result.ExecArgs {
				widgets = append(widgets, g.Label(
					fmt.Sprintf("  '%s'", arg),
				))
			}
			widgets = append(widgets, g.Label("]"))
			widgets = append(widgets, g.Label(""))
		}

		/* --------------------- SYSTEMD-RUN -------------------- */
		systemdRunOutput := util.ParseSystemdRunOutput(myCmd.Result.Output)

		widgets = append(widgets, g.Line(
			g.ArrowButton("Arrow", g.DirectionRight),
			g.Label("SYSTEMD-RUN"),
		))

		for _, keyValue := range systemdRunOutput {
			key := keyValue[0]
			value := keyValue[1]

			widgets = append(widgets, g.Line(
				g.Button(key).OnClick(func() {
					util.CopyToClipboard(value)
				}),
				g.Label(value),
			))
		}

		widgets = append(widgets, g.Label(""))

		/* ------------------------- KEY ------------------------ */
		widgets = append(widgets, g.Line(
			g.ArrowButton("Arrow", g.DirectionRight),
			g.Label("KEY"),
		))

		widgets = append(widgets, g.Line(
			g.Button("As").OnClick(func() {
				util.CopyToClipboard(myCmd.Keybind.As)
			}),
			g.Label(myCmd.Keybind.As),
		))

		widgets = append(widgets, g.Line(
			g.Button("Cmd").OnClick(func() {
				util.CopyToClipboard(myCmd.Keybind.Cmd)
			}),
			g.Label(myCmd.Keybind.Cmd),
		))

		widgets = append(widgets, g.Line(
			g.Button("Args").OnClick(func() {
				util.CopyToClipboard(strings.Join(myCmd.Keybind.Args, " "))
			}),
			g.Label("["),
		))
		for _, arg := range myCmd.Keybind.Args {
			widgets = append(widgets, g.Label(
				fmt.Sprintf("  '%s'", arg),
			))
		}
		widgets = append(widgets, g.Label("]"))

		widgets = append(widgets, g.Line(
			g.Button("Wait").OnClick(func() {
				util.CopyToClipboard(fmt.Sprintf("%t", myCmd.Keybind.Wait))
			}),
			g.Label(fmt.Sprintf("%t", myCmd.Keybind.Wait)),
		))

		widgets = append(widgets, g.Line(
			g.Button("Key").OnClick(func() {
				util.CopyToClipboard(myCmd.KeybindKey)
			}),
			g.Label(myCmd.KeybindKey),
		))

		widgets = append(widgets, g.Line(
			g.Button("Mod").OnClick(func() {
				util.CopyToClipboard(myCmd.KeybindMod)
			}),
			g.Label(myCmd.KeybindMod),
		))

		widgets = append(widgets, g.Label(""))

		/* --------------------- RAW OUTPUT --------------------- */
		if myCmd.Result.Output != "" {
			widgets = append(widgets, g.Line(
				g.ArrowButton("Arrow", g.DirectionRight),
				g.Label("RAW OUTPUT"),
			))
			widgets = append(widgets, g.Button("Copy Raw Output").OnClick(func() {
				util.CopyToClipboard(myCmd.Result.Output)
			}))
			widgets = append(widgets, g.Label(myCmd.Result.Output))
		}
	}

	g.SingleWindow("Runner").Layout(
		widgets...,
	)
}
