package service

import "fmt"

// BizError 业务错误
type BizError struct {
	Code   string
	ErrMsg string
}

func (receiver BizError) Error() string {
	return fmt.Sprintf("code:%s,errmsg:%s", receiver.Code, receiver.ErrMsg)
}

func NewBizErr(code, msg string) *BizError {
	return &BizError{
		Code:   code,
		ErrMsg: msg,
	}
}
