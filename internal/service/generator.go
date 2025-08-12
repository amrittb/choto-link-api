package service

// "time"

// var seed uint64 = uint64(time.Now().UnixMicro())
var seed uint64 = 10_000_000_000_000

// var seed uint64 = uint64(1)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Create a map for fast lookup of character values
var base62Map = func() map[rune]uint64 {
	m := make(map[rune]uint64)
	for i, c := range base62Chars {
		m[c] = uint64(i)
	}
	return m
}()

func EncodeBase62(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []byte
	for num > 0 {
		remainder := num % 62
		encoded = append([]byte{base62Chars[remainder]}, encoded...)
		num /= 62
	}
	return string(encoded)
}

func GetNextId() uint64 {
	seed += 1
	return seed
}
