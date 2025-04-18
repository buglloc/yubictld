package yubictl

import "time"

type AcquireRsp struct {
	ID     string `json:"id"`
	Serial uint32 `json:"serial"`
}

type PingReq struct {
	ID string `json:"id"`
}

type RebootReq struct {
	ID string `json:"id"`
}

type TouchReq struct {
	ID       string        `json:"id"`
	Delay    time.Duration `json:"delay"`
	Duration time.Duration `json:"duration"`
}

type ReleaseReq struct {
	ID string `json:"id"`
}
