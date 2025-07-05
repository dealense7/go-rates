package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	dto "github.com/dealense7/go-rate-app/internal/DTO"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/eko/gocache/store/redis/v4"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strconv"
	"time"
)

type CacheStoreRepository struct {
	tag   string
	db    *sqlx.DB
	cache *redis.RedisStore
	repo  interfaces.StoreRepository
}

func NewCacheStoreRepository(db *sqlx.DB, cache *redis.RedisStore) interfaces.StoreRepository {
	return &CacheStoreRepository{
		tag:   "products",
		db:    db,
		cache: cache,
		repo:  NewMySQLStoreRepository(db),
	}
}
func (c CacheStoreRepository) GetProductById(id int) (models.SingleProductItem, error) {
	return c.repo.GetProductById(id)
}

func (c CacheStoreRepository) GetForSlider(ctx context.Context) ([]models.ProductItem, error) {
	var items []models.ProductItem

	// You may think cache is useless here, but it is not if there are many visitors
	// Generate cache key
	cacheKey := c.generateCacheKey("products-main-slider")

	// Check if data exists in Cache
	if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
		if jsonStr, ok := cached.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
				return items, nil
			}
		}
	}

	// Get data from Database if not found in Cache
	items, err := c.repo.GetForSlider(ctx)
	if err != nil {
		return nil, err
	}

	// Cache data
	if data, err := json.Marshal(items); err == nil {
		_ = c.cache.Set(ctx, cacheKey, data, store.WithExpiration(time.Second*20))
	}

	return items, nil
}

func (c CacheStoreRepository) GetItemsList(ctx context.Context, page int, category int, name string) ([]models.ProductItem, error) {
	var items []models.ProductItem

	// Generate Cache Key
	cacheKey := c.generateCacheKey("products-list-page:%d-category:%d", page, category)

	// Temporary: don't use cache for filter with name, I have no limits on redis
	if name != "" {
		return c.repo.GetItemsList(ctx, page, category, name)
	}

	// Return data from Cache if found
	if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
		if jsonStr, ok := cached.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
				return items, nil
			}
		}
	}

	// Fetch items from Database if not found in Cache
	items, err := c.repo.GetItemsList(ctx, page, category, name)
	if err != nil {
		return nil, err
	}

	// JSON encode and cache fetched data
	if data, err := json.Marshal(items); err == nil {
		_ = c.cache.Set(ctx, cacheKey, data)
	}

	return items, nil
}

func (c CacheStoreRepository) GetItemsCount(ctx context.Context, category int, name string) (int, error) {
	var count int

	// Generate Cache Key
	cacheKey := c.generateCacheKey("products-list-items-count-category:%d", category)

	// Temporary: If user filters with name don't use cache
	if name != "" {
		return c.repo.GetItemsCount(ctx, category, name)
	}

	// Get if data is cached
	if val, err := c.cache.Get(ctx, cacheKey); err == nil {
		if countStr, ok := val.(string); ok {
			count, err := strconv.Atoi(countStr)
			if err == nil {
				return count, nil
			}
		}
	}

	// calculate count again if not in cache
	count, err := c.repo.GetItemsCount(ctx, category, name)
	if err != nil {
		return 0, err
	}

	// Set fetched data in cache
	if err := c.cache.Set(ctx, cacheKey, count); err != nil {
		fmt.Println(err.Error())
	}

	return count, nil
}

func (c CacheStoreRepository) GetForCategorySlider(ctx context.Context) ([]models.CategorySlider, error) {
	var items []models.CategorySlider

	// Generate Cache Key
	cacheKey := c.generateCacheKey("products-main-category-slider")

	// Check if data is already cached
	if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
		if jsonStr, ok := cached.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
				return items, nil
			}
		}
	}

	// Fetch items if not cached
	items, err := c.repo.GetForCategorySlider(ctx)
	if err != nil {
		return nil, err
	}

	// JSON encode and cache
	if data, err := json.Marshal(items); err == nil {
		if err := c.cache.Set(ctx, cacheKey, data, store.WithExpiration(time.Minute*5)); err != nil {
			slog.Error("failed to set cache", "key", cacheKey, "error", err)
		}
	}

	return items, nil
}

func (c CacheStoreRepository) GetProductByBarCode(barCode string) (int64, error) {
	return c.repo.GetProductByBarCode(barCode)
}

func (c CacheStoreRepository) CreateItem(data dto.Product) (int64, error) {
	return c.repo.CreateItem(data)
}

func (c CacheStoreRepository) AddOrUpdatePrice(itemId int64, data dto.ProductPrice) error {
	return c.repo.AddOrUpdatePrice(itemId, data)
}

func (c CacheStoreRepository) DisableOldPrices() error {
	return c.repo.DisableOldPrices()
}

func (c CacheStoreRepository) generateCacheKey(suffix string, args ...interface{}) string {
	base := fmt.Sprintf("tag:%s:%s", c.tag, suffix)
	if len(args) > 0 {
		return fmt.Sprintf(base, args...)
	}
	return base
}
