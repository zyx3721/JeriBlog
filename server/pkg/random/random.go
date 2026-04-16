/*
项目名称：JeriBlog
文件名称：random.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：随机数生成工具
*/

package random

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
const digits = "0123456789"
const alphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ExistsChecker 检查字符串是否存在的函数类型
type ExistsChecker func(string) (bool, error)

// String 生成指定长度的随机字符串（大小写字母+数字）
func String(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = alphaNum[b[i]%byte(len(alphaNum))]
	}
	return string(b)
}

// Code 生成指定长度的随机码（小写字母+数字）
func Code(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// Digits 生成指定长度的纯数字随机码（用于验证码等场景）
func Digits(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = digits[b[i]%10]
	}
	return string(b)
}

// UniqueCode 生成唯一的随机码（最多尝试100次）
func UniqueCode(length int, check ExistsChecker) (string, error) {
	for i := 0; i < 100; i++ {
		code := Code(length)
		exists, err := check(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return Code(length), nil
}
