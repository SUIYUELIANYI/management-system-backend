package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response请求响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (r *Response) Error() string {
	return r.Msg
}

func SendResponse(w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	if err == nil {
		httpx.WriteJson(w, http.StatusOK, Success(resp))
	} else {
		if e, ok := err.(*Response); ok {
			httpx.WriteJson(w, http.StatusOK, e)
		} else {
			httpx.WriteJson(w, http.StatusBadRequest, struct {
				Message string `json:"message"`
			}{Message: err.Error()})
		}
	}
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Data: data,
		Msg:  "", // 如果不想要Msg字段在返回时出现可以直接用nil，下面的Data同理
	}
}

func Error(code int, message string) *Response {
	return &Response{
		Code: code,
		Data: "",
		Msg:  message,
	}
}
