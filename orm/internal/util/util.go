package util

import (
	"regexp"
	"strings"
)

var (
	reCamel = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// chatgpt写的
// 问题:  你能用go代码通过正则的方式, 把驼峰式字符串转为下划线分割的字符串吗
func CamelToUnderline(s string) string {
	// 将驼峰式字符串中的大写字母和小写字母+数字的分界点用 "_" 分隔开来
	converted := reCamel.ReplaceAllStringFunc(s, func(s string) string {
		return strings.Join([]string{s[:1], "_", s[1:]}, "")
	})
	// 将字符串转为小写
	converted = strings.ToLower(converted)
	return converted
}
