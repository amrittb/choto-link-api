package repository

var db map[string]string = make(map[string]string)

func Save(mapping UrlMap) {
	db[mapping.ShortUrl] = mapping.LongUrl
}

func Get(shortUrl string) (string, bool) {
	val, ok := db[shortUrl]
	return val, ok
}
