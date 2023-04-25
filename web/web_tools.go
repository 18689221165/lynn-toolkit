package web

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// UIntParam 得到Uint64类型的参数
func UIntParam(c *gin.Context, paramName string) (uint64, error) {
	oidstr := c.Param(paramName)

	// 当路由为/xxxx/*param时，param是可选的，取的oid为"/param"
	// 若访问/xxxx时，会重定向到/xxxx/;这里取的oid为"/"
	oidstr = strings.TrimPrefix(oidstr, "/")
	if oidstr == "" {
		oidstr = c.Query(paramName)
	}
	if oidstr == "" {
		oidstr = c.PostForm(paramName)
	}

	return strconv.ParseUint(oidstr, 10, 64)
}

func GetJwtToken(c *gin.Context, authKey string) string {
	if r := c.GetHeader(authKey); r != "" {
		return r
	}
	if r, ok := c.GetQuery(authKey); ok {
		return r
	}
	if r, ok := c.GetPostForm(authKey); ok {
		return r
	}
	if r, err := c.Cookie(authKey); err == nil {
		return r
	}
	return ""
}
