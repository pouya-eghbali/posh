package utils

func Unquote(s string) string {
	return s[1 : len(s)-1]
}
