package validates

import "goim/server"

// 过滤敏感词
func VerificationFilter(data string) bool {
	b, _ := server.Filter.FindIn(data)
	return b
}
