package xerr

// OK 成功返回
const OK uint32 = 200
const InternalError uint32 = 500

/**(前3位代表业务,后三位代表具体功能)**/

// ServerCommonError 全局错误码
const ServerError uint32 = 200001
const ParamError uint32 = 200002
const ForbiddenError uint32 = 200003
const DbError uint32 = 200005
const CacheError = 200006
const NotifyError uint32 = 200007
const GetGlobalConfigError uint32 = 200008
const RecordNotFound = 200009
const DbRecordExist = 200014
const UnmarshalFailed = 200018
const TokenInvalid = 200019 // 内部鉴权错误
const QueryError = 200020
const UpdateError = 200021
const CreateError = 200022

// 用户
const (
	// UserNotExists         uint32 = 200101
	UserPwdInvalid        uint32 = 200102
	UserTokenInvalid      uint32 = 200103
	UserLoginExpired      uint32 = 200104
	UserHasNoPermission   uint32 = 200105
	UserDisabled          uint32 = 200106 // 用户已被禁用
	UserExpired           uint32 = 200107
	DownloadLimitExceeded uint32 = 200108 // 下载次数超限
	UserNotLogin          uint32 = 200109 // 用户未登录
)

// 商品
const (
	SkuNotExists uint32 = 200201 // 商品不存在
	SkuNotUp     uint32 = 200202 // 商品未上架
	SkuDown      uint32 = 200203 // 商品已下线
)

// 退款
const (
	RefundError uint32 = 200301 // 退款业务错误
)
