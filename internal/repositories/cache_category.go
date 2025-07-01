package repositories

import (
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
	"github.com/eko/gocache/store/redis/v4"
	"github.com/jmoiron/sqlx"
)

type CacheCategoryRepository struct {
	db    *sqlx.DB
	cache *redis.RedisStore
	repo  interfaces.CategoryRepository
}

func NewCacheCategoryRepository(db *sqlx.DB, cache *redis.RedisStore) interfaces.CategoryRepository {
	return &CacheCategoryRepository{
		db:    db,
		cache: cache,
		repo:  NewMySQLCategoryRepository(db),
	}
}

func (r *CacheCategoryRepository) GetItems() ([]models.Category, error) {
	return r.repo.GetItems()
}
