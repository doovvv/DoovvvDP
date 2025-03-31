package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomCode(length int) string {
	// 字符集（数字）
	charset := []rune("0123456789")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 生成随机验证码
	code := make([]rune, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
func RandomString(length int) string {
	// 字符集（数字和字母）
	charset := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// 生成随机验证码
	code := make([]rune, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}