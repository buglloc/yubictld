package config

import (
	"fmt"

	"github.com/buglloc/yubictld/internal/httpd"
)

type ServerCfg struct {
	Addr string `koanf:"addr"`
}

func (r *Runtime) NewServer() (*httpd.Server, error) {
	yk, err := r.YkMan()
	if err != nil {
		return nil, fmt.Errorf("create ykman runtime: %w", err)
	}

	touch, err := r.Toucher()
	if err != nil {
		return nil, fmt.Errorf("create touch runtime: %w", err)
	}

	return httpd.NewServer(
		httpd.WithAddr(r.cfg.Server.Addr),
		httpd.WithYkMan(yk),
		httpd.WithToucher(touch),
	)
}
