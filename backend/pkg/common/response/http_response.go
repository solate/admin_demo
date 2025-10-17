package response

import (
	"context"
	"net/http"

	"admin_backend/pkg/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type Response struct {
	Code int    `json:"code"` // 状态码
	Msg  string `json:"msg"`  // 提示信息
	Data any    `json:"data"`
}

func Success(data any) *Response {
	return &Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: data,
	}
}

func Fail(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}

func OkHanandler(_ context.Context, v any) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {

		//错误返回
		errcode := xerr.ServerError
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}

		logx.WithContext(ctx).Errorf("【API-ERR】 : %+v ", err)

		return http.StatusOK, Fail(int(errcode), errmsg) // 还是返回200， 错误码和错误信息
	}
}
