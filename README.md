# cactus

A hotkey launcher

## Motivation

I created this because xbindkeys, sxhkd, dxhkd, etc. didn't launch new processes within the context of the current user's systemd default cgroup slice.

This causes a number of side effects, such as the following appearing in my journal after restarting

```text
sxhkd.service: Found left-over process 1182850 (bash) in control group while starting unit. Ignoring.
This usually indicates unclean termination of a previous run, or service implementation deficiencies.
```

The output of all running commands would output to `sxhkd`'s controlling terminal; since I personally start `sxhkd` as a service, the output shows up in the unit journal.

Not only that, but applications that quit unsuccessfully or dumped their cores stoped the sxhkd service, consequently necessitating a full restart of `sxhkd.service`. Since my application launcher (`dmenu`/`rofi`) is started by by `sxhkd`, I am only able to restart if I already have a terminal open or by switching to a different virtual console

`Cactus` solves this by decoupling execution context by adding an indireciton layer who's environment is controlled by systemd

## Features

- Choose shell or exec per command to use
- View list of all available commands so you don't have to remember it right away
- By default, pass arguments directly to exec-style function (faster startup)
- Execs out to `systemd-run` behind the scenes
- Config hot reload (you may have to hover your mouse over interface for an update due to ImgGUI)
- Escape to close window

## Usage

- AUR PKGBUILD TODO

```sh
git clone https://github.com/eankeen/cactus
cd cactus
go install .
```

You can use `cactus` via the terminal, _or_ launch via a traditional hotkey daemon

Ex. `sxhkdrc`

```txt
super + y
	cactus
```

## Configuration

```toml
# $XDG_CONFIG_HOME/cactus/binds.toml

# If `run` is `exec` or is absent, `cactus` executes directly to
# e.g. `execv` using `cmd` and `args`
[B]
run = "exec"
cmd = "brave"

[Control-B]
cmd = "brave"

[Shift-B]
cmd = "vlc"
args = ["--rate", "0.5", "--play-and-stop"]

[Alt-B]
cmd = "vlc"
args = ["--loop", "--ignore-filetypes"]


# Else if `run` is `dash`, or `bash`, `cactus` executes using
# e.g. `bash` using `cmd` only
[Shift-J]
run = "sh"
cmd = "cd ~/repos/sticker-selector && go run ."

[J]
run = "bash"
cmd = "</dev/null ls -al
```

Find the full keymap list [here](./keymap/keymap.go)

## Troubleshooting

You may have issues with launching applications because the systemd slice may not match your current environment. Ammend your `~/.profile` etc.

```sh
systemctl --user import-environment LANG LANGUAGE LC_ALL PATH XDG_CONFIG_HOME XDG_DATA_HOME
```

Remember, if you are launching `cactus` through another hotkey daemon that is launched as a systemd user service, you will have to restart it after this

## TODO

- showSuccess (show the output on success)
- gorouting when execing only when wait = true (systemd-run returns immediately)
- delete systemd service unit after wars (or have a respective setting)
- ad hoc keybindings to change auxillary settings
- other TODOs
- fix horizontal scrollbar

## Potential Improvements

- add font
- support --tty and stream or stream from journal
