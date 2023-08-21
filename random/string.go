package random

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// String generates a random string of the specified length using the default charset.
func String(length int) string {
	return stringWithCharset(length, charset)
}

func stringWithCharset(length int, charset string) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		chi := defaultRnd.Intn(len(charset))
		b[i] = charset[chi]
	}
	return string(b)
}
