package dto

type ResponseData struct {
	Code    int         `json:"code"`    // Ma status code
	Message string      `json:"message"` // Thong bao loi
	Data    interface{} `json:"data"`    // Du lieu duoc return
}

type ErrResponseData struct {
	Code   int         `json:"code"`   // Ma status code
	Error  string      `json:"error"`  //
	Detail interface{} `json:"detail"` // Thong bao loi
}
