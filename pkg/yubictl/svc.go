package yubictl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type SvcClient struct {
	pingInterval time.Duration
	httpc        *resty.Client
}

func NewSvcClient(upstream string, opts ...Option) *SvcClient {
	c := &SvcClient{
		pingInterval: DefaultPingInterval,
		httpc: resty.New().
			SetJSONEscapeHTML(false).
			SetHeader("Content-Type", "application/json").
			SetBaseURL(upstream).
			SetRetryCount(3).
			SetRetryWaitTime(1 * time.Second).
			SetRetryMaxWaitTime(10 * time.Second).
			AddRetryCondition(func(rsp *resty.Response, err error) bool {
				return err != nil || rsp.StatusCode() == http.StatusInternalServerError
			}),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *SvcClient) Acquire(ctx context.Context) (*Yubikey, error) {
	var out AcquireRsp
	var serviceErr ServiceError
	rsp, err := c.httpc.R().
		SetContext(ctx).
		SetError(&serviceErr).
		SetResult(&out).
		ForceContentType("application/json").
		Post("/v1/acquire")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !rsp.IsSuccess() {
		if serviceErr.Code != ServiceErrorCodeNone {
			return nil, &serviceErr
		}

		return nil, fmt.Errorf("request failed: non-200 status code: %s", rsp.Status())
	}

	if out.ID == "" || out.Serial == 0 {
		return nil, fmt.Errorf("server returns unexpected response: %s", rsp.String())
	}

	yCtx, yCancel := context.WithCancel(context.Background())
	yk := &Yubikey{
		id:        out.ID,
		serial:    out.Serial,
		httpc:     c.httpc,
		pingTick:  c.pingInterval,
		ctx:       yCtx,
		cancelCtx: yCancel,
		closed:    make(chan struct{}),
	}
	go yk.pingLoop()

	return yk, nil
}
