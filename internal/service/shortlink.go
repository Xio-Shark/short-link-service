package service

import (
	"context"
	urlpkg "net/url"
	"time"
)

type LinkRepository interface {
	Create(ctx context.Context, url string, expireAt int64) (string, error)
	Resolve(ctx context.Context, code string, access AccessInfo) (string, error)
	Stats(ctx context.Context, code string) (int64, int64, error)
}

type ShortLinkService struct {
	Repo LinkRepository
}

type AccessInfo struct {
	IP        string
	UserAgent string
}

func (s *ShortLinkService) Create(ctx context.Context, url string, expireAt string) (string, error) {
	if url == "" {
		return "", ErrInvalidRequest
	}
	if _, err := urlpkg.ParseRequestURI(url); err != nil {
		return "", ErrInvalidRequest
	}

	var expireAtUnix int64
	if expireAt != "" {
		t, err := time.Parse(time.RFC3339, expireAt)
		if err != nil {
			return "", ErrInvalidRequest
		}
		expireAtUnix = t.Unix()
	}

	return s.Repo.Create(ctx, url, expireAtUnix)
}

func (s *ShortLinkService) Resolve(ctx context.Context, code string, access AccessInfo) (string, error) {
	return s.Repo.Resolve(ctx, code, access)
}

func (s *ShortLinkService) Stats(ctx context.Context, code string) (int64, int64, error) {
	return s.Repo.Stats(ctx, code)
}
