package ykman

import "time"

const (
	DefaultLockTTL = time.Hour
)

type Option func(*YkMan)

func WithLockTTL(ttl time.Duration) Option {
	return func(y *YkMan) {
		y.lockTTL = ttl
	}
}

func WithDiscovery(discovery Discovery) Option {
	return func(y *YkMan) {
		y.discovery = discovery
	}
}
