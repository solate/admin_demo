package constants

// 默认密码
const UserPwdDefault = "123456"

// 角色
const (
	RoleAdminRoot = 1 // 超级管理员
	RoleAgency    = 2 // 代理商
	RoleStore     = 3 // 商家
)

// 性别
const (
	GenderMan   = 1 // 男
	GenderWoman = 2 // 女
)

// SmsCodeValidityTime 短信有效期 毫秒单位
const SmsCodeValidityTime = 5 * 60 * 1000
