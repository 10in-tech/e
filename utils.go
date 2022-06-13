package e

import "bytes"

// 拼接字符串
func join(str ...string) string {
	var buffer bytes.Buffer
	for _, s := range str {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func inArray(str string, array []string) bool {
	for _, i := range array {
		if i == str {
			return true
		}
	}
	return false
}
