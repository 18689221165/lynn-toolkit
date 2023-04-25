package service

var (
	ErrInternalServerError = NewApiError("PUB_INTERANL_ERR", "服务器内部错误")         // 公共错误-服务器内部错误
	ErrNotFound            = NewApiError("PUB_NOT_FOUND_ERR", "资源未找到")          // 公共错误-记录不存在
	ErrConflict            = NewApiError("PUB_RECORD_CONFLICT", "记录已存在")        // 公共错误-记录已存在
	ErrBadParamInput       = NewApiError("PUB_BAD_PARAM", "参数错误")               // 公共错误-参数错误
	ErrHttpMethod          = NewApiError("PUB_HTTP_METHOD_ERR", "不支持的HTTP请求方法") // 公共错误-不支持的HTTP请求方法
)
