package base62

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(num uint64) string {
	if num == 0 {
		return "0"
	}

	encoded := make([]byte, 0)
	for num > 0 {
		remainder := num % 62
		encoded = append([]byte{base62Chars[remainder]}, encoded...)
		num /= 62
	}
	return string(encoded)
}
