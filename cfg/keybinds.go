package cfg

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type Keybinds map[string]KeybindEntry
type KeybindEntry struct {
	As   string   `toml: "run"`
	Cmd  string   `toml: "cmd"`
	Args []string `toml: "args"`

	// Auxilary members
	Wait            bool `toml: "wait"`
	OutputOnSuccess bool `toml: "outputOnSuccess"`
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
	keybindings, err := km.parseConfig()
	if err != nil {
		return err
	}
	km.keybinds = keybindings

	return nil
}

func (kbm *keybindsMnger) parseConfig() (Keybinds, error) {
	cfgText, err := ioutil.ReadFile(kbm.path)
	if err != nil {
		return Keybinds{}, fmt.Errorf("Error: Could not read the file '%s'. Ensure the file exists and has the proper permissions\n%w", kbm.path, err)
	}

	var keyBinds Keybinds
	err = toml.Unmarshal(cfgText, &keyBinds)
	if err != nil {
		return Keybinds{}, fmt.Errorf("Error: Could not unmarshal contents of file '%s'. Ensure the file is valid TOML\n%w", kbm.path, err)
	}

	return keyBinds, nil
}
