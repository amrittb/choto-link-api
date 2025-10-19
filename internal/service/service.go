package service

import (
	"github.com/amrittb/choto-link-api/internal/base62"
	"github.com/amrittb/choto-link-api/internal/generator"
	"github.com/amrittb/choto-link-api/internal/repository"
)

func GenerateShortUrl(longUrl string) string {
	// 1. Get nextId from current Id Range
	// 2. Encode nextId to base62
	// 3. Save shortUrl into DB
	// 5. Return response
	shortUrl := base62.Encode(generator.GetNextId())
	repository.Save(repository.UrlMap{ShortUrl: shortUrl, LongUrl: longUrl})
	return shortUrl
}

func GetLongUrl(shortUrl string) (string, bool) {
	return repository.Get(shortUrl)
}
