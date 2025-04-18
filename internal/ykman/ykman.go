package ykman

import (
	"fmt"
	"sync"
	"time"

	"github.com/buglloc/fidoctl"
	"github.com/rs/zerolog/log"
)

type YkMan struct {
	lockTTL   time.Duration
	discovery Discovery
	mu        sync.Mutex
	store     []*Yubikey
}

func NewYkMan(opts ...Option) *YkMan {
	yk := &YkMan{
		lockTTL: DefaultLockTTL,
	}

	for _, opt := range opts {
		opt(yk)
	}

	return yk
}

func (y *YkMan) ReloadDevices() error {
	y.mu.Lock()
	defer y.mu.Unlock()

	devices, err := fidoctl.Enumerate()
	if err != nil {
		return fmt.Errorf("enumerate devices: %w", err)
	}

	y.store = y.store[:0]
	for _, dev := range devices {
		yk, err := newYubikey(dev, y.discovery)
		if err != nil {
			return fmt.Errorf("create yubikey %s: %w", dev.String(), err)
		}

		if y.discovery != nil && yk.Port() == 0 {
			continue
		}

		y.store = append(y.store, yk)
	}

	return nil
}

func (y *YkMan) Acquire(clientID string) (*Yubikey, error) {
	y.mu.Lock()
	defer y.mu.Unlock()

	for _, yk := range y.store {
		if !yk.IsFree() {
			if time.Since(yk.lastAccess) < y.lockTTL {
				continue
			}

			log.Warn().
				Uint32("yk_serial", yk.serial).
				Time("last_access", yk.lastAccess).
				Dur("stale_time", time.Since(yk.lastAccess)).
				Msg("release stale lock")

			if err := yk.Release(); err != nil {
				log.Error().
					Uint32("yk_serial", yk.serial).
					Time("last_access", yk.lastAccess).
					Msg("release failed")
				continue
			}
		}

		err := yk.Acquire(clientID)
		if err != nil {
			continue
		}

		return yk, nil
	}

	return nil, ErrNoFreeYubikey
}

func (y *YkMan) ForClient(clientID string) (*Yubikey, error) {
	y.mu.Lock()
	defer y.mu.Unlock()

	for _, yk := range y.store {
		if yk.client != clientID {
			continue
		}

		return yk, nil
	}

	return nil, fmt.Errorf("get Yubikey for client %s: %w", clientID, ErrNoAssociated)
}

func (y *YkMan) Devices() []*Yubikey {
	y.mu.Lock()
	defer y.mu.Unlock()

	out := make([]*Yubikey, len(y.store))
	copy(out, y.store)
	return out
}
