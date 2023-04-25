package service

var (
	FAIL = NewApiResult("fail", "失败", nil) // 操作失败
)

func Succ(data interface{}) *ApiResult {
	return NewApiResult("success", "成功", data)
}
