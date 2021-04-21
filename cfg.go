package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type CfgToml map[string]CfgEntry
type CfgEntry struct {
	Cmd string `toml: "cmd"`
	Run string `toml: "run"`
}

func getCfgFile() string {
	var cfgDir string

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		cfgDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(cfgDir, "cactus", "binds.toml")
}

func getCfg(cfgFile string) CfgToml {
	cfgText, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Config file '%s' does not exist\n", cfgFile)
		os.Exit(1)
	}
	handle(err)

	var cfgToml CfgToml
	err = toml.Unmarshal(cfgText, &cfgToml)
	handle(err)

	return cfgToml
}
