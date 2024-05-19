package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/elanq/pastebin-go/model"
	"github.com/elanq/pastebin-go/repository"
	"github.com/go-redis/redis/v8"
)

var (
	UserIDNotFoundError = errors.New("User ID not found")
	EmptyURLError       = errors.New("URL should not be empty")
	CacheTypeError      = errors.New("Error while getting value from cache")
)

const (
	shortURLCharNum      = 10
	defaultCacheDuration = 10 * time.Minute

	urlCacheKeyPrefix = "url_cache_key"
)

type URL interface {
	Create(context.Context, model.URL) (*model.URL, error)
	Update(context.Context, model.URL) error
	Delete(context.Context, string) error
	GetByShortURL(context.Context, string) (*model.URL, error)
}

type url struct {
	urlRepository repository.URL
	cacheService  Cache
}

// Create implements URL.
func (u *url) Create(ctx context.Context, req model.URL) (*model.URL, error) {
	err := validateCreate(req)
	if err != nil {
		return nil, err
	}
	generateShortURL(&req)

	createdURL, err := u.urlRepository.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	cacheKey := urlCacheKey(req.ShortURL)
	err = u.cacheService.Set(ctx, cacheKey, *createdURL, defaultCacheDuration)
	if err != nil {
		return nil, err
	}
	return createdURL, nil
}

// Delete implements URL.
func (u *url) Delete(context.Context, string) error {
	panic("unimplemented")
}

// GetByShortURL implements URL.
func (u *url) GetByShortURL(ctx context.Context, url string) (*model.URL, error) {
	cacheKey := urlCacheKey(url)
	res, err := u.cacheService.Get(ctx, cacheKey)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if res != nil {
		s, ok := res.(string)
		if !ok {
			return nil, CacheTypeError
		}
		var urlResult model.URL
		err = json.Unmarshal([]byte(s), &urlResult)
		return &urlResult, err
	}

	urlModel, err := u.urlRepository.FindByShortUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	u.cacheService.Set(ctx, cacheKey, urlModel, defaultCacheDuration)
	return urlModel, nil
}

// Update implements URL.
func (u *url) Update(context.Context, model.URL) error {
	panic("unimplemented")
}

func urlCacheKey(shortURL string) string {
	return strings.Join([]string{urlCacheKeyPrefix, shortURL}, ":")
}

func validateCreate(req model.URL) error {
	if req.UserId == "" {
		return UserIDNotFoundError
	}
	if req.LongURL == "" {
		return EmptyURLError
	}
	return nil
}

func generateShortURL(req *model.URL) {
	shortURL := uniuri.NewLen(shortURLCharNum)
	req.ShortURL = shortURL
}

func NewURL(urlRepository repository.URL, cacheService Cache) URL {
	return &url{
		urlRepository: urlRepository,
		cacheService:  cacheService,
	}
}
