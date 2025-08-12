package service

var db map[string]string = make(map[string]string)

func GenerateShortUrl(request string) string {
	shortUrl := EncodeBase62(GetNextId())
	db[shortUrl] = request
	return shortUrl
}

func GetLongUrl(shortUrl string) (string, bool) {
	val, ok := db[shortUrl]
	return val, ok
}
