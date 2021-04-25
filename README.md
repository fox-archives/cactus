# cactus

Hotkey launcher

## Description

Meant to be used as as a secondary application launcher. As an alternative a separate launch mode you can create
with i3, etc.

NOTE

Doesn't work as it's intended purpuse, fails when not ran within a tty

Ex. sxhkdrc

```txt
super + y
	cactus
```

## Configuration

```toml
# $XDG_CONFIG_HOME/cactus/binds.toml
[A]
cmd = "pavucontrol"

[O]
cmd = "obs"

[C]
cmd = "code"
```

## Key Names

From glfw

https://pkg.go.dev/github.com/AllenDang/giu@v0.5.3?utm_source=gopls#Key

## TOD

- proper error handling
