package util

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
