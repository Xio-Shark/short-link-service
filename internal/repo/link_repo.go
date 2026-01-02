package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"short-link-service/internal/model"
	"short-link-service/internal/service"
)

const maxCodeRetry = 5

type LinkRepo struct {
	DB    *gorm.DB
	Cache *redis.Client
}

func NewLinkRepo(db *gorm.DB, cache *redis.Client) *LinkRepo {
	return &LinkRepo{DB: db, Cache: cache}
}

func (r *LinkRepo) Create(ctx context.Context, url string, expireAt int64) (string, error) {
	now := time.Now().Unix()
	link := &model.Link{
		OriginalURL: url,
		ExpireAt:    expireAt,
		Status:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tx := r.DB.WithContext(ctx).Begin()
	if err := tx.Create(link).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	code := service.EncodeBase62(link.ID)
	var lastErr error
	for i := 0; i < maxCodeRetry; i++ {
		err := tx.Model(&model.Link{}).
			Where("id = ?", link.ID).
			Update("code", code).Error
		if err == nil {
			link.Code = code
			lastErr = nil
			break
		}
		if isDuplicateKey(err) {
			code = code + service.EncodeBase62(uint64(time.Now().UnixNano()%62))
			lastErr = err
			continue
		}
		tx.Rollback()
		return "", err
	}
	if lastErr != nil {
		tx.Rollback()
		return "", lastErr
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	r.cacheSet(ctx, link.Code, link.OriginalURL, link.ExpireAt)
	return link.Code, nil
}

func (r *LinkRepo) Resolve(ctx context.Context, code string, access service.AccessInfo) (string, error) {
	if r.Cache != nil {
		if val, err := r.Cache.Get(ctx, code).Result(); err == nil && val != "" {
			return val, nil
		}
	}

	var link model.Link
	now := time.Now().Unix()
	err := r.DB.WithContext(ctx).
		Where("code = ? AND status = 1 AND (expire_at = 0 OR expire_at > ?)", code, now).
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", service.ErrNotFound
		}
		return "", err
	}

	r.cacheSet(ctx, link.Code, link.OriginalURL, link.ExpireAt)
	r.writeVisit(ctx, link.ID, access)
	return link.OriginalURL, nil
}

func (r *LinkRepo) Stats(ctx context.Context, code string) (int64, int64, error) {
	var link model.Link
	now := time.Now().Unix()
	err := r.DB.WithContext(ctx).
		Where("code = ? AND status = 1 AND (expire_at = 0 OR expire_at > ?)", code, now).
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, service.ErrNotFound
		}
		return 0, 0, err
	}

	var pv int64
	if err := r.DB.WithContext(ctx).Model(&model.Visit{}).
		Where("link_id = ?", link.ID).
		Count(&pv).Error; err != nil {
		return 0, 0, err
	}

	var uv int64
	if err := r.DB.WithContext(ctx).Model(&model.Visit{}).
		Where("link_id = ?", link.ID).
		Distinct("ip").
		Count(&uv).Error; err != nil {
		return 0, 0, err
	}

	return pv, uv, nil
}

func (r *LinkRepo) cacheSet(ctx context.Context, code string, url string, expireAt int64) {
	if r.Cache == nil {
		return
	}
	var ttl time.Duration
	if expireAt > 0 {
		ttl = time.Until(time.Unix(expireAt, 0))
		if ttl <= 0 {
			return
		}
	}
	_ = r.Cache.Set(ctx, code, url, ttl).Err()
}

func (r *LinkRepo) writeVisit(ctx context.Context, linkID uint64, access service.AccessInfo) {
	if access.IP == "" && access.UserAgent == "" {
		return
	}
	visit := &model.Visit{
		LinkID:    linkID,
		IP:        access.IP,
		UserAgent: access.UserAgent,
		CreatedAt: time.Now().Unix(),
	}
	_ = r.DB.WithContext(ctx).Create(visit).Error
}

func isDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
