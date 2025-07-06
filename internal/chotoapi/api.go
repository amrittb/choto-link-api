package chotoapi

import (
	"crypto/md5"
	"encoding/hex"
)

var db map[string]string = make(map[string]string)

// Takes in long url and returns long url
func Create(request string) string {
	hasher := md5.New()
	hasher.Write([]byte(request))
	shortUrl := hex.EncodeToString(hasher.Sum(nil))
	db[shortUrl] = request
	return shortUrl
}

func Get(shortUrl string) (string, bool) {
	val, ok := db[shortUrl]
	return val, ok
}
