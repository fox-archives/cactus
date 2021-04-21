package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type CfgToml map[string]struct {
	Cmd string `toml: "cmd"`
	Run string `toml: "run"`
}

func getCfgDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.Getenv("XDG_CONFIG_HOME")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".config")
	}
}

func getCfg() CfgToml {
	// TODO: use ini
	cfgFile := filepath.Join(getCfgDir(), "cactus", "bindings.toml")
	var cfgToml CfgToml

	cfgText, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Config file '%s' does not exist\n", cfgFile)
		os.Exit(1)
	}
	handle(err)

	err = toml.Unmarshal(cfgText, &cfgToml)
	handle(err)

	return cfgToml
}
