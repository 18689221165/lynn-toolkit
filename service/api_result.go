package service

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ApiResult 微服务公有的结果
type ApiResult struct {
	code string
	msg  string
	data interface{}
}

func NewApiResult(code, msg string, data interface{}) *ApiResult {
	return &ApiResult{code, msg, data}
}

func NewApiError(code, msg string) *ApiResult {
	return &ApiResult{code: code, msg: msg}
}

func (res ApiResult) Error() string {
	return fmt.Sprintf("code=%s,msg=%s", res.code, res.msg)
}

func (res *ApiResult) Code(code string) *ApiResult {
	res.code = code
	return res
}

func (res *ApiResult) Msg(msg string) *ApiResult {
	res.msg = msg
	return res
}

func (res *ApiResult) Data(data interface{}) *ApiResult {
	res.data = data
	return res
}

func (res ApiResult) MarshalJSON() ([]byte, error) {
	builder := &strings.Builder{}
	builder.WriteString("{\"code\":\"")
	builder.WriteString(res.code)
	builder.WriteString("\"")

	builder.WriteString(",\"msg\":\"")
	builder.WriteString(res.msg)
	builder.WriteString("\"")

	if res.data != nil {
		bytes, _ := json.Marshal(res.data)
		builder.WriteString(",\"data\":")
		builder.Write(bytes)
	}

	builder.WriteString("}")

	return []byte(builder.String()), nil
}
