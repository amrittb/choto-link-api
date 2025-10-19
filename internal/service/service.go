package service

import (
	"github.com/amrittb/choto-link-api/internal/allocator"
	"github.com/amrittb/choto-link-api/internal/base62"
	"github.com/amrittb/choto-link-api/internal/repository"
)

type ShortUrlService struct {
	idAllocator   *allocator.IdAllocator
	urlRepository *repository.UrlRepository
}

func NewShortUrlService(idAllocator *allocator.IdAllocator, urlRepository *repository.UrlRepository) *ShortUrlService {
	return &ShortUrlService{idAllocator: idAllocator, urlRepository: urlRepository}
}

func (service *ShortUrlService) AllocateAndSaveShortUrl(longUrl string) (string, error) {
	// 1. Get nextId from current Id Range
	// 2. Encode nextId to base62
	// 3. Save shortUrl into DB
	// 5. Return response
	nextId, err := service.idAllocator.NextId()
	if err != nil {
		return "", err
	}

	shortUrl := base62.Encode(nextId)
	err = service.urlRepository.Save(repository.UrlMap{ShortUrl: shortUrl, LongUrl: longUrl})
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

func (service *ShortUrlService) GetLongUrl(shortUrl string) (string, bool) {
	return service.urlRepository.Get(shortUrl)
}
