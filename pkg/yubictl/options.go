package yubictl

import "time"

const (
	DefaultPingInterval = 5 * time.Second
)

type Option func(*SvcClient)

func WithPingInterval(d time.Duration) Option {
	return func(c *SvcClient) {
		c.pingInterval = d
	}
}

type TouchOption func(r *TouchReq)

func TouchWithDuration(d time.Duration) TouchOption {
	return func(r *TouchReq) {
		r.Duration = d
	}
}

func TouchWithDelay(d time.Duration) TouchOption {
	return func(r *TouchReq) {
		r.Delay = d
	}
}
