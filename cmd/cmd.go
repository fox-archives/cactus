package run

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/eankeen/cactus/cfg"
	"github.com/google/uuid"
)

type Cmd struct {
	KeybindKey string
	KeybindMod string
	Keybind    cfg.KeybindEntry
	HasRan     bool
	Result     CmdResult
}

type CmdResult struct {
	ExecName string
	ExecArgs []string
	Err      error
	Output   string
}

func New() *Cmd {
	// This sets defaults
	return &Cmd{
		KeybindKey: "",
		KeybindMod: "",
		// Not defaults, overriden in RunCmdOnce
		Keybind: cfg.KeybindEntry{
			Cmd:  "",
			As:   "",
			Wait: false,
		},
		HasRan: false,
		Result: CmdResult{},
	}
}

func (cmd *Cmd) RunCmdOnce(mod string, key string, keybindEntry cfg.KeybindEntry) {
	// If the result is not nil, we already ran a command
	if cmd.HasRan == true {
		return
	}

	cmd.KeybindMod = mod
	cmd.KeybindKey = key
	cmd.Keybind = keybindEntry
	cmd.HasRan = true

	// If runCmd() fails to properly run command, it stops execution on it's own
	cmd.Result = cmd.runCmd()

	if cmd.HasRan && !cmd.Keybind.InfoOnSuccess {
		os.Exit(0)
	}
}

func (cmd *Cmd) runCmd() CmdResult {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return CmdResult{
			ExecName: "",
			ExecArgs: []string{},
			Err:      fmt.Errorf("Cactus Internal Error: Could not generate random number\n%w", err),
			Output:   "",
		}
	}

	args := []string{
		"--no-ask-password",
		"--unit", "cactus-" + uuid.String(),
		"--description", fmt.Sprintf("Cactus Start for command: '%s'", cmd.Keybind.Cmd),
		"--send-sighup",
		"--working-directory", os.Getenv("HOME"),
		"--user",
	}

	if cmd.Keybind.Wait {
		args = append(args, "--wait")
	}

	args = append(args, "--")

	// TODO custom ones
	switch cmd.Keybind.As {
	case "sh":
		_, err := os.Stat("/usr/bin/dash")
		if errors.Is(err, os.ErrNotExist) {
			args = append(args, "/usr/bin/sh", "-c", cmd.Keybind.Cmd)
		} else {
			// use dash directly if possible, even 'sh' is specified since
			// sometimes 'sh' is symlinked to bash
			args = append(args, "/usr/bin/dash", "-c", cmd.Keybind.Cmd)
		}
	case "bash":
		args = append(args, "/usr/bin/bash", "-c", cmd.Keybind.Cmd)
	default:
		cmd.Keybind.As = "exec"
		args = append(args, cmd.Keybind.Cmd)
		args = append(args, cmd.Keybind.Args...)
	}

	execName := "/usr/bin/systemd-run"
	execCmd := exec.Command(execName, args...)

	output, err := execCmd.CombinedOutput()
	return CmdResult{
		ExecName: execName,
		ExecArgs: args,
		Err:      err,
		Output:   string(output),
	}

}
