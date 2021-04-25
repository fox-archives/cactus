package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type CfgCactus struct {
	// Full path to font
	FontFile string
	// Size of the font
	FontSize int
}

type CfgBind map[string]CfgEntry
type CfgEntry struct {
	Cmd string `toml: "cmd"`
	Run string `toml: "run"`
}

func getCfgFile(file string) string {
	var cfgDir string

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		cfgDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(cfgDir, "cactus", file)
}

func getCfgBinds(cfgFile string) CfgBind {
	cfgText, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Config file '%s' does not exist\n", cfgFile)
		os.Exit(1)
	}
	handle(err)

	var cfgBind CfgBind
	err = toml.Unmarshal(cfgText, &cfgBind)
	handle(err)

	return cfgBind
}

func getCfgCactus(cfgFile string) CfgCactus {
	cfgText, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Config file '%s' does not exist\n", cfgFile)
		os.Exit(1)
	}
	handle(err)

	var cfgCactus CfgCactus
	err = toml.Unmarshal(cfgText, &cfgCactus)
	handle(err)

	if cfgCactus.FontSize == 0 {
		cfgCactus.FontSize = 16
	}

	cfgCactus.FontFile = os.ExpandEnv(cfgCactus.FontFile)

	return cfgCactus
}
