package cfg

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type Keybinds map[string]KeybindEntry
type KeybindEntry struct {
	Cmd string `toml: "cmd"`
	Run string `toml: "run"`
}

type keybindsMnger struct {
	// Path to the keybindings
	path string

	// Keybinds config map
	keybinds Keybinds
}

func NewKeybindsMnger(keybindsFile string) keybindsMnger {
	return keybindsMnger{
		path:     keybindsFile,
		keybinds: Keybinds{},
	}
}

func (km *keybindsMnger) Get() *Keybinds {
	return &km.keybinds
}

func (km *keybindsMnger) Reload() error {
	keybindings, err := km.parseConfig(km.path)
	if err != nil {
		return err
	}
	km.keybinds = keybindings

	return nil
}

// TODO bindingsFile repetitive
func (km *keybindsMnger) parseConfig(bindingsFile string) (Keybinds, error) {
	cfgText, err := ioutil.ReadFile(bindingsFile)
	if err != nil {
		return Keybinds{}, fmt.Errorf("Error: Could not read the file '%s'. Ensure the file exists and has the proper permissions\n%w", bindingsFile, err)
	}

	var keyBinds Keybinds
	err = toml.Unmarshal(cfgText, &keyBinds)
	if err != nil {
		return Keybinds{}, fmt.Errorf("Error: Could not unmarshal contents of file '%s'. Ensure the file is valid TOML\n%w", bindingsFile, err)
	}

	return keyBinds, nil
}
