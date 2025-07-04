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
	cacheKey := fmt.Sprintf("tag:%s:products-main-slider", c.tag)
	var items []models.ProductItem

	if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
		if jsonStr, ok := cached.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
				return items, nil
			}
		}
	}

	items, err := c.repo.GetForSlider(ctx)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(items); err == nil {
		_ = c.cache.Set(ctx, cacheKey, data, store.WithExpiration(time.Second*20))
	}

	return items, nil
}

func (c CacheStoreRepository) GetItemsList(ctx context.Context, page int, category int, name string) ([]models.ProductItem, error) {
	cacheKey := fmt.Sprintf("tag:%s:products-list-page:%d-category:%d", c.tag, page, category)
	var items []models.ProductItem

	// I don't need many keys in redis, I pay 5$ for vps :)))
	if name != "" {
		return c.repo.GetItemsList(ctx, page, category, name)
	}

	if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
		if jsonStr, ok := cached.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
				return items, nil
			}
		}
	}

	items, err := c.repo.GetItemsList(ctx, page, category, name)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(items); err == nil {
		_ = c.cache.Set(ctx, cacheKey, data)
	}

	return items, nil
}

func (c CacheStoreRepository) GetItemsCount(ctx context.Context, category int, name string) (int, error) {
	cacheKey := fmt.Sprintf("tag:%s:products-list-items-count-category:%d", c.tag, category)

	var count int
	if name != "" {
		return c.repo.GetItemsCount(ctx, category, name)
	}

	if val, err := c.cache.Get(ctx, cacheKey); err == nil {
		if countStr, ok := val.(string); ok {
			count, err := strconv.Atoi(countStr)
			if err == nil {
				return count, nil
			}
		}
	}

	count, err := c.repo.GetItemsCount(ctx, category, name)
	if err != nil {
		return 0, err
	}

	if err := c.cache.Set(ctx, cacheKey, count); err != nil {
		fmt.Println(err.Error())
	}

	return count, nil
}

func (c CacheStoreRepository) GetForCategorySlider(ctx context.Context) ([]models.CategorySlider, error) {
	cacheKey := fmt.Sprintf("tag:%s:products-main-category-slider", c.tag)
	var items []models.CategorySlider

	if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
		if jsonStr, ok := cached.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
				return items, nil
			}
		}
	}

	items, err := c.repo.GetForCategorySlider(ctx)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(items); err == nil {
		_ = c.cache.Set(ctx, cacheKey, data, store.WithExpiration(time.Minute*5))
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
