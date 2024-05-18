package common

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func NewResponse(code int, msg string, data any) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewSuccessResponse(data any) *Response {
	return NewResponse(0, "success", data)
}

func NewErrorResponse(code int, msg string) *Response {
	return NewResponse(code, msg, nil)
}
