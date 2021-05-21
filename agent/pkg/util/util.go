package util

import "encoding/json"

const (
	Base64 string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
)

func Contains(arr []string, element string) bool {
	for _, item := range arr {
		if item == element {
			return true
		}
	}
	return false
}

func Convert(from interface{}, to interface{}) {
	rs, _ := json.Marshal(from)
	_ = json.Unmarshal(rs, to)
}

func MaskLeft(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = '*'
	}
	return string(rs)
}
