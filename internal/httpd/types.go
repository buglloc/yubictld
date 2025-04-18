package httpd

type Config struct {
	Addr string
}

type ErrorRsp struct {
	Code    int    `json:"err_code"`
	Message string `json:"err_message"`
}
