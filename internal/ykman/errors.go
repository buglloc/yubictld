package ykman

import "errors"

var ErrNoFreeYubikey = errors.New("no free Yubikey was found")
var ErrNoAssociated = errors.New("associated Yubikey not found")
