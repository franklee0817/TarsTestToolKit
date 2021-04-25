package tools

import (
	"math/rand"
)

var alphas = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

// GetStr 获取定长字符串
func GetStr(len uint32) string {
	s := ""
	for i := uint32(0); i < len; i++ {
		s += alphas[rand.Intn(36)]
	}

	return s
}
