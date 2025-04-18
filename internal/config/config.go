package config

import (
	"fmt"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/buglloc/yubictld/internal/httpd"
	"github.com/buglloc/yubictld/internal/touchctl"
	"github.com/buglloc/yubictld/internal/ykman"
)

type Config struct {
	Server ServerCfg `koanf:"server"`
	Touch  TouchCfg  `koanf:"touch"`
	YkMan  YkManCfg  `koanf:"ykman"`
}

func (c *Config) Validate() error {
	return nil
}

type Runtime struct {
	cfg     *Config
	toucher touchctl.Toucher
	ykman   *ykman.YkMan
}

func LoadConfig(files ...string) (*Config, error) {
	out := &Config{
		Server: ServerCfg{
			Addr: httpd.DefaultAddr,
		},
		Touch: TouchCfg{
			Kind: touchctl.ToucherKindH4ptix,
		},
		YkMan: YkManCfg{
			LockTTL:   time.Hour,
			Discovery: ykman.DiscoveryKindToucher,
		},
	}

	k := koanf.New(".")
	if err := k.Load(env.Provider("YUBICTL", "_", nil), nil); err != nil {
		return nil, fmt.Errorf("load env config: %w", err)
	}

	yamlParser := yaml.Parser()
	for _, fpath := range files {
		if err := k.Load(file.Provider(fpath), yamlParser); err != nil {
			return nil, fmt.Errorf("load %q config: %w", fpath, err)
		}
	}

	return out, k.Unmarshal("", &out)
}

func (c *Config) NewRuntime() (*Runtime, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &Runtime{
		cfg: c,
	}, nil
}
