package repository

type UrlRepository struct {
	db map[string]string
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{db: make(map[string]string)}
}

func (repo *UrlRepository) Save(mapping UrlMap) error {
	repo.db[mapping.ShortUrl] = mapping.LongUrl
	return nil
}

func (repo *UrlRepository) Get(shortUrl string) (string, bool) {
	val, ok := repo.db[shortUrl]
	return val, ok
}
