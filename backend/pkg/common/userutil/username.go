package userutil

import (
	"fmt"
	"regexp"
	"strings"
)

// GenerateDefaultUsername 生成默认用户名
// 格式: "课研所用户" + 手机号后四位
// 如果手机号为空或长度不足4位，则返回 "课研所用户"
func GenerateDefaultUsername(phone string) string {
	// 移除空格和特殊字符
	phone = strings.TrimSpace(phone)

	// 如果手机号为空或长度不足4位，返回默认名称
	if len(phone) < 4 {
		return "课研所用户"
	}

	// 获取手机号后四位
	lastFourDigits := phone[len(phone)-4:]

	// 生成带后缀的用户名
	return fmt.Sprintf("课研所用户%s", lastFourDigits)
}

// NormalizePhone 规范化手机号码
// 移除国际区号前缀(+86)、空格、横线等特殊字符，保留纯数字手机号
// 例如: +86-18270548941 -> 18270548941
//
//	+86 18270548941 -> 18270548941
//	+8618270548941 -> 18270548941
func NormalizePhone(phone string) string {
	if phone == "" {
		return ""
	}

	// 移除空格
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}

	// 移除常见的国际区号前缀和分隔符
	// 支持格式: +86-xxx, +86 xxx, +86xxx
	re := regexp.MustCompile(`^\+86[-\s]?`)
	phone = re.ReplaceAllString(phone, "")

	// 移除其他分隔符：横线、空格等
	re = regexp.MustCompile(`[-\s]`)
	phone = re.ReplaceAllString(phone, "")

	// 只保留数字
	re = regexp.MustCompile(`\D`)
	phone = re.ReplaceAllString(phone, "")

	return phone
}

// MaskPhone 对手机号进行脱敏处理
// 将中间四位替换为 ****
// 例如: 17612345678 -> 176****5678
func MaskPhone(phone string) string {
	// 先规范化手机号
	phone = NormalizePhone(phone)

	// 如果手机号为空或长度不足11位，直接返回原值
	if len(phone) < 11 {
		return phone
	}

	// 取前3位和后4位，中间用****替换
	return phone[:3] + "****" + phone[len(phone)-4:]
}
