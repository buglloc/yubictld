package touchctl

import (
	"fmt"
	"time"

	"github.com/buglloc/h4ptix/software/h4ptix"
)

var _ Toucher = (*H4ptix)(nil)

type H4ptix struct {
	h *h4ptix.H4ptix
}

func NewH4ptix(opts ...H4ptixOption) (*H4ptix, error) {
	var h4ptixOpts []h4ptix.Option
	for _, opt := range opts {
		switch v := opt.(type) {
		case optH4ptixSerial:
			if v.serial != "" {
				h4ptixOpts = append(h4ptixOpts, h4ptix.WithDeviceSerial(v.serial))
			}

		default:
			return nil, fmt.Errorf("invalid h4ptix option: %T", opt)
		}
	}

	h, err := h4ptix.NewH4ptix(h4ptixOpts...)
	if err != nil {
		return nil, fmt.Errorf("initialize H4ptix: %w", err)
	}

	return &H4ptix{
		h: h,
	}, nil
}

func (h *H4ptix) Location() string {
	return h.h.Location()
}

func (h *H4ptix) Touch(port int, delay time.Duration, duration time.Duration) error {
	return h.h.Trigger(h4ptix.TriggerReq{
		Port:     port,
		Duration: duration,
		Delay:    delay,
	})
}

func (h *H4ptix) Close() error {
	return h.h.Close()
}
