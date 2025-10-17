package xerr

import (
	"fmt"
)

type CodeError struct {
	errCode uint32
	errMsg  string
}

// GetErrCode 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// GetErrMsg 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}
func NewErrCode(errCode uint32) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ServerError, errMsg: errMsg}
}

// IsTokenError 判断是否是token相关的错误
func IsTokenError(err error) bool {
	if codeErr, ok := err.(*CodeError); ok {
		// token相关的错误码
		switch codeErr.errCode {
		case TokenInvalid, UserTokenInvalid, UserLoginExpired, UserNotLogin:
			return true
		}
	}
	return false
}
