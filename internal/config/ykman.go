package config

import (
	"fmt"
	"time"

	"github.com/buglloc/yubictld/internal/ykman"
)

type YkManCfg struct {
	LockTTL   time.Duration       `koanf:"lock_ttl"`
	Discovery ykman.DiscoveryKind `koanf:"discovery"`
	Manual    struct {
		Yubikeys []struct {
			Serial uint32 `koanf:"serial"`
			Port   int    `koanf:"port"`
		} `koanf:"yubikeys"`
	} `koanf:"manual"`
}

func (r *Runtime) YkMan() (*ykman.YkMan, error) {
	if r.ykman != nil {
		return r.ykman, nil
	}

	disco, err := r.NewDiscovery()
	if err != nil {
		return nil, fmt.Errorf("inialize yubikeys discovery: %w", err)
	}

	yk := ykman.NewYkMan(
		ykman.WithLockTTL(r.cfg.YkMan.LockTTL),
		ykman.WithDiscovery(disco),
	)
	return yk, yk.ReloadDevices()
}

func (r *Runtime) NewDiscovery() (ykman.Discovery, error) {
	switch r.cfg.YkMan.Discovery {
	case ykman.DiscoveryKindNone:
		return nil, nil

	case ykman.DiscoveryKindManual:
		yMap := make(map[uint32]int)
		for _, y := range r.cfg.YkMan.Manual.Yubikeys {
			yMap[y.Serial] = y.Port
		}

		return ykman.NewManualDiscovery(yMap)

	case ykman.DiscoveryKindToucher:
		toucher, err := r.Toucher()
		if err != nil {
			return nil, fmt.Errorf("initialize toucher: %w", err)
		}

		return ykman.NewToucherDiscovery(toucher)

	default:
		return nil, fmt.Errorf("unsupported discovery: %s", r.cfg.YkMan.Discovery)
	}
}
