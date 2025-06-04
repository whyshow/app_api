package utils

import "regexp"

// ValidateEmailFormat 验证邮箱格式
func ValidateEmailFormat(email string) bool {
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
