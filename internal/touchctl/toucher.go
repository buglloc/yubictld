package touchctl

import (
	"errors"
	"time"
)

type Toucher interface {
	Touch(port int, delay time.Duration, duration time.Duration) error
	Location() string
}

type NopToucher struct{}

func NewNopToucher() *NopToucher {
	return &NopToucher{}
}

func (t *NopToucher) Touch(_ int, _ time.Duration, _ time.Duration) error {
	return errors.New("not implemented")
}

func (t *NopToucher) Location() string {
	return ""
}
