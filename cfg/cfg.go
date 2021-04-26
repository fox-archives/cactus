package cfg

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

type Cfg struct {
	// Full path to font
	FontFile string

	// Size of the font
	FontSize int
}

type cfgMnger struct {
	// Path to the config file
	path string

	// Config struct
	cfg Cfg
}

func NewCfgMnger(cfgFile string) cfgMnger {
	return cfgMnger{
		path: cfgFile,
		cfg:  Cfg{},
	}
}

func (c *cfgMnger) Get() *Cfg {
	return &c.cfg
}
func (c *cfgMnger) Reload() error {
	cfg, err := c.parseConfig(c.path)
	if err != nil {
		return err
	}
	c.cfg = cfg

	return nil
}

func (c *cfgMnger) parseConfig(cfgFile string) (Cfg, error) {
	cfgText, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return Cfg{}, fmt.Errorf("Error: Could not read the file '%s'. Ensure the file exists and has the proper permissions\n%w", cfgFile, err)
	}

	var cfg Cfg
	err = toml.Unmarshal(cfgText, &cfg)
	if err != nil {
		return Cfg{}, fmt.Errorf("Error: Could not unmarshal contents of file '%s'. Ensure the file is valid TOML\n%w", cfgFile, err)
	}

	// defaults
	cfg.FontFile = os.ExpandEnv(cfg.FontFile)
	if cfg.FontSize == 0 {
		cfg.FontSize = 16
	}

	return cfg, nil
}
