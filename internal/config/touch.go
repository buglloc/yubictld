package config

import (
	"fmt"

	"github.com/buglloc/yubictld/internal/touchctl"
)

type TouchCfg struct {
	Kind   touchctl.ToucherKind `koanf:"kind"`
	H4ptix struct {
		Serial string `koanf:"serial"`
	} `koanf:"h4ptix"`
}

func (r *Runtime) Toucher() (touchctl.Toucher, error) {
	if r.toucher != nil {
		return r.toucher, nil
	}

	switch r.cfg.Touch.Kind {
	case touchctl.ToucherKindNone:
		r.toucher = touchctl.NewNopToucher()
		return r.toucher, nil

	case touchctl.ToucherKindH4ptix:
		var err error
		r.toucher, err = touchctl.NewH4ptix(
			touchctl.H4ptixWithSerial(r.cfg.Touch.H4ptix.Serial),
		)
		return r.toucher, err

	default:
		return nil, fmt.Errorf("unknown touch kind %s", r.cfg.Touch.Kind)
	}
}
