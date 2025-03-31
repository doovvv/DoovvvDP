package utils

import (
	"regexp"
)

// 检验手机号
func IsValidPhoneNumber(phone string) bool {
	// 正则表达式：以1开头，第二位是3-9中的一个，后面是9个数字
	// 这个正则表达式适用于中国的手机号
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}