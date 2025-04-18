package ykman

import (
	"encoding"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/buglloc/yubictld/internal/touchctl"
)

var _ encoding.TextUnmarshaler = (*DiscoveryKind)(nil)
var _ encoding.TextMarshaler = (*DiscoveryKind)(nil)

type DiscoveryKind string

const (
	DiscoveryKindNone    DiscoveryKind = ""
	DiscoveryKindToucher DiscoveryKind = "toucher"
	DiscoveryKindManual  DiscoveryKind = "manual"
)

func (k *DiscoveryKind) UnmarshalText(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "", "none":
		*k = DiscoveryKindNone
	case "toucher":
		*k = DiscoveryKindToucher
	case "manual":
		*k = DiscoveryKindManual
	default:
		return fmt.Errorf("invalid discovery kind: %s", string(data))
	}
	return nil
}

func (k DiscoveryKind) MarshalText() ([]byte, error) {
	return []byte(k), nil
}

type Discovery interface {
	Port(y *Yubikey) int
}

type ManualDiscovery struct {
	yubikeys map[uint32]int
}

func NewManualDiscovery(yubikeys map[uint32]int) (*ManualDiscovery, error) {
	return &ManualDiscovery{
		yubikeys: yubikeys,
	}, nil
}

func (d *ManualDiscovery) Port(y *Yubikey) int {
	return d.yubikeys[y.Serial()]
}

type ToucherDiscovery struct {
	touch touchctl.Toucher
}

func NewToucherDiscovery(touch touchctl.Toucher) (*ToucherDiscovery, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("toucher discovery is no supported on %s so far", runtime.GOOS)
	}

	return &ToucherDiscovery{
		touch: touch,
	}, nil
}

func (d *ToucherDiscovery) Port(y *Yubikey) int {
	yubiLocation := y.Location()
	touchLocation := d.touch.Location()
	if yubiLocation == "" || touchLocation == "" {
		return 0
	}

	if d.hubLocation(yubiLocation) != d.hubLocation(touchLocation) {
		return 0
	}

	return d.portLocation(yubiLocation)
}

func (d *ToucherDiscovery) hubLocation(loc string) string {
	idx := strings.LastIndexByte(loc, '.')
	if idx == -1 {
		return ""
	}

	return loc[:idx]
}

func (d *ToucherDiscovery) portLocation(loc string) int {
	idx := strings.LastIndexByte(loc, '.')
	if idx == -1 {
		return 0
	}

	port, _ := strconv.Atoi(loc[idx+1:])
	return port
}
