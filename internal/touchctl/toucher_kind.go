package touchctl

import (
	"encoding"
	"fmt"
	"strings"
)

var _ encoding.TextUnmarshaler = (*ToucherKind)(nil)
var _ encoding.TextMarshaler = (*ToucherKind)(nil)

type ToucherKind string

const (
	ToucherKindNone   ToucherKind = ""
	ToucherKindH4ptix ToucherKind = "h4ptix"
)

func (k *ToucherKind) UnmarshalText(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "", "none":
		*k = ToucherKindNone
	case "h4ptix":
		*k = ToucherKindH4ptix
	default:
		return fmt.Errorf("invalid toucher kind: %s", string(data))
	}
	return nil
}

func (k ToucherKind) MarshalText() ([]byte, error) {
	return []byte(k), nil
}
