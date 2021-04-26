# cactus

A hotkey launcher

## Motivation

I created this because xbindkeys, sxhkd, dxhkd, etc. didn't launch new processes within the context of the current user's systemd default cgroup slice.

This causes a number of side effects, such as the following appearing in my journal after restarting

```text
sxhkd.service: Found left-over process 1182850 (bash) in control group while starting unit. Ignoring.
This usually indicates unclean termination of a previous run, or service implementation deficiencies.
```

The output of all running commands would show in the journal and applications that quit unsuccessfully stoped the sxhkd service, consequently necessitating a full restart of `sxhkd.service`

## Features

- Choose shell per command to use
- By default, pass arguments directly to exec-style function (faster startup)
- Shells out to `systemd-run` behind the scenes
- Config hot reload (you may have to hover your mouse over interface for an update due to ImgGUI)

## Setup

Ex. sxhkdrc

```txt
super + y
	cactus
```

## Usage

- AUR PKGBUILD TODO

```sh
git clone https://github.com/eankeen/cactus
cd cactus
go install .
```

## Configuration

```toml
# $XDG_CONFIG_HOME/cactus/binds.toml
[A]
cmd = "pavucontrol"

[O]
cmd = "obs"

[S]
run = "dash"
cmd = "cd ~/repos/sticker-selector && go run ."
```

## Key Names

From glfw

[key.go](./key.go)

https://pkg.go.dev/github.com/AllenDang/giu@v0.5.3?utm_source=gopls#Key

## TODO

- automatick binds to show journal / output

## Potential Improvements

- add font
- use init format
