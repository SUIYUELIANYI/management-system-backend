package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 改了几个版本的错误处理，决定还是用下面这个版本，looklook的太复杂了，下面这个能满足所有返回（无论对错）都是Response形式，也方便对接前端
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
