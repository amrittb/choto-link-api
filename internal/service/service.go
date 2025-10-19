package service

import "github.com/amrittb/choto-link-api/internal/base62"

var db map[string]string = make(map[string]string)

func GenerateShortUrl(request string) string {
	shortUrl := base62.EncodeBase62(GetNextId())
	db[shortUrl] = request
	return shortUrl
}

func GetLongUrl(shortUrl string) (string, bool) {
	val, ok := db[shortUrl]
	return val, ok
}
