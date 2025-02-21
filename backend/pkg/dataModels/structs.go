package datamodels

type Response struct {
	Code     int    `json:"code"`
	Status   string `json:"status"`
	ErrorMsg string `json:"error_msg"`
}
