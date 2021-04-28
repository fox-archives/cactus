package run

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/eankeen/cactus/cfg"
	"github.com/eankeen/cactus/util"
	"github.com/google/uuid"
)

type Cmd struct {
	KeybindKey string
	KeybindMod string
	Keybind    cfg.KeybindEntry
	HasRan     bool
	Result     *CmdResult
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
		Keybind: cfg.KeybindEntry{
			Cmd:  "",
			Run:  "",
			Wait: false,
		},
		HasRan: false,
		Result: nil,
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
}

func (cmd *Cmd) runCmd() *CmdResult {
	uuid, err := uuid.NewRandom()
	if err != nil {
		util.Handle(fmt.Errorf("Error: Could not generate random number\n%w", err))
	}

	args := []string{
		"--no-ask-password",
		"--unit", "cactus-" + uuid.String(),
		"--description", fmt.Sprintf("Cactus Start for command: '%s'", cmd.Keybind.Cmd),
		"--send-sighup",
		"--working-directory",

		os.Getenv("HOME"), "--user",
	}

	if cmd.Keybind.Wait {
		args = append(args, "--wait")
	}

	switch cmd.Keybind.Run {
	case "dash":
		args = append(args, "/usr/bin/dash", "-c", cmd.Keybind.Cmd)
	case "bash":
		args = append(args, "/usr/bin/bash", "-c", cmd.Keybind.Cmd)
	default:
		args = append(args, cmd.Keybind.Cmd)
	}

	execName := "/usr/bin/systemd-run"
	execCmd := exec.Command(execName, args...)

	output, err := execCmd.CombinedOutput()
	return &CmdResult{
		ExecName: execName,
		ExecArgs: args,
		Err:      err,
		Output:   string(output),
	}

}
