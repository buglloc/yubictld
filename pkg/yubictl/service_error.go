package yubictl

import "fmt"

type ServiceErrorCode int

const (
	ServiceErrorCodeNone ServiceErrorCode = iota
	ServiceErrorInternalError
	ServiceErrorNoFreeYubikey
)

type ServiceError struct {
	HttpCode int              `json:"-"`
	Code     ServiceErrorCode `json:"error_code"`
	Msg      string           `json:"message"`
}

func (e *ServiceError) Is(err error) bool {
	other, ok := err.(*ServiceError)
	if !ok {
		return false
	}

	if e == nil && other == nil {
		return true
	}
	if e == nil || other == nil {
		return false
	}
	return e.Code == other.Code
}

func (e *ServiceError) IsPermanent() bool {
	if e.Code == ServiceErrorCodeNone {
		return false
	}

	if e.Code == ServiceErrorInternalError {
		return false
	}

	return true
}

func (e *ServiceError) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("yubictl_svc (%d): %s", e.Code, e.Msg)
	}
	return fmt.Sprintf("yubictl_svc (%d)", e.Code)
}
