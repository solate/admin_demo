package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	message[ServerError] = "服务器开小差啦,稍后再来试一试"
	message[ParamError] = "参数错误"
	message[ForbiddenError] = "无权限操作"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[NotifyError] = "!!! 请联系开发人员确认操作 !!!"
	message[GetGlobalConfigError] = "获取全局配置失败"
	message[RecordNotFound] = "未找到该记录"
	message[InternalError] = "内部错误"

	message[DbRecordExist] = "数据库记录已存在"
	message[UnmarshalFailed] = "数据序列化错误"
	message[QueryError] = "查询错误"
	message[UpdateError] = "更新错误"
	message[CreateError] = "保存错误"

	// message[UserNotExists] = "用户不存在"
	message[UserPwdInvalid] = "用户密码错误"
	message[UserTokenInvalid] = "无效的token"
	message[UserLoginExpired] = "用户登录已过期"
	message[UserHasNoPermission] = "用户无权限"
	message[TokenInvalid] = "鉴权错误"

	message[UserDisabled] = "用户已被禁用"
	message[UserExpired] = "用户已过期"
	message[DownloadLimitExceeded] = "下载次数超限"
	message[UserNotLogin] = "用户未登录"

}

func MapErrMsg(errCode uint32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errCode uint32) bool {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}

func NewErrCodeMsgByCode(errCode uint32) *CodeError {
	return &CodeError{
		errCode: errCode,
		errMsg:  MapErrMsg(errCode),
	}
}
