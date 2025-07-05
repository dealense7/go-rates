package repositories

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	dto "github.com/dealense7/go-rate-app/internal/DTO"
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"time"
)

type MySQLStoreRepository struct {
	db *sqlx.DB
}

func NewMySQLStoreRepository(db *sqlx.DB) interfaces.StoreRepository {
	return &MySQLStoreRepository{db: db}
}

func (r *MySQLStoreRepository) GetProductById(id int) (models.SingleProductItem, error) {
	var item models.SingleProductItem

	query := `SELECT sp.id,
				   sp.name_ka AS name,
				   sp.company,
				   sp.image_url,
				   JSON_ARRAYAGG(
						   JSON_OBJECT(
								   'id', spp.id,
								   'date', DATE_FORMAT(spp.created_at, '%Y-%m-%d %H:%i:%s'),
								   'price', spp.price,
								   'provider', s.name,
								   'provider_logo', s.logo_url
						   )
				   ) AS prices
			FROM store_products AS sp
					 JOIN store_product_prices spp ON sp.id = spp.product_id and spp.status = 1
					 JOIN store_providers s ON spp.store_id = s.id
			WHERE sp.id = ?
			GROUP BY sp.id, sp.name_ka, sp.company, sp.image_url`

	if err := r.db.Get(&item, query, id); err != nil {
		return models.SingleProductItem{}, fmt.Errorf("failed to get product %d: %w", id, err)
	}

	return item, nil
}

func (r *MySQLStoreRepository) GetForSlider(ctx context.Context) ([]models.ProductItem, error) {
	var items []models.ProductItem

	const query = `
		WITH interesting_products AS (
			SELECT 
				spp.product_id,
				MAX(spp.price) - MIN(spp.price) AS price_diff
			FROM store_product_prices spp
			JOIN store_products sp 
				ON sp.id = spp.product_id 
				AND sp.status = true
			WHERE spp.status = true
			GROUP BY spp.product_id
			HAVING COUNT(*) > 1
			ORDER BY price_diff DESC
			LIMIT 50
		), random_products AS (
			SELECT product_id
			FROM interesting_products
			ORDER BY RAND()
			LIMIT 6
		)
		SELECT 
			sp.id,
			sp.name_ka AS name,
			sp.company,
			sp.image_url,
			sp.volume,
			sp.origin,
			MIN(spp.price) AS min_price,
			MAX(spp.price) AS max_price
		FROM store_products sp
		JOIN random_products rp 
			ON rp.product_id = sp.id
		JOIN store_product_prices spp 
			ON spp.product_id = sp.id 
			AND spp.status = true
		GROUP BY 
			sp.id, 
			sp.name_ka, 
			sp.company, 
			sp.image_url, 
			sp.volume, 
			sp.origin
		LIMIT 6`

	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("failed to get slider products: %w", err)
	}

	return items, nil
}
func (r *MySQLStoreRepository) GetItemsList(ctx context.Context, page int, category int, name string) ([]models.ProductItem, error) {
	var items []models.ProductItem
	offset := 30

	// Base query
	query := `SELECT sp.id,
                     sp.name_ka     AS name,
                     sp.company,
                     sp.image_url,
                     sp.volume,
                     sp.origin,
                     MIN(spp.price) AS min_price,
                     MAX(spp.price) AS max_price
              FROM store_products AS sp
                       JOIN store_product_prices AS spp
                            ON spp.product_id = sp.id AND spp.status = true`

	// Parameters for the query
	params := []interface{}{}

	// Add category join and filter if category > 1
	if category > 1 {
		query += ` JOIN golang.categories c ON sp.category_id = c.id
                   WHERE sp.status = true AND c.parent_id = ?`
		params = append(params, category)
	} else {
		query += ` WHERE sp.status = true`
	}

	// Add name filter for name_ka OR name_en if provided also dont make it to long
	if name != "" && len(name) < 255 {
		query += ` AND (LOWER(sp.name_ka) LIKE LOWER(?) OR LOWER(sp.name_en) LIKE LOWER(?) OR LOWER(sp.company) LIKE LOWER(?))`
		params = append(params, "%"+name+"%", "%"+name+"%", "%"+name+"%")
	}

	// Complete the query
	query += ` GROUP BY sp.id, sp.name_ka, sp.company, sp.image_url, sp.volume, sp.origin
               HAVING COUNT(spp.product_id) > 1
               ORDER BY (max_price - min_price) DESC LIMIT 30 OFFSET ?`

	// Append offset parameter
	params = append(params, offset*page)

	// Execute query
	err := r.db.Select(&items, query, params...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *MySQLStoreRepository) GetItemsCount(ctx context.Context, category int, name string) (int, error) {
	var count int

	// Base query
	query := `SELECT COUNT(*)
              FROM (
                  SELECT spp.product_id
                  FROM store_product_prices AS spp
                           JOIN store_products sp ON sp.id = spp.product_id`

	// Parameters for the query
	params := []interface{}{}

	// Add category join and filter if category > 1
	if category > 1 {
		query += ` JOIN golang.categories c ON sp.category_id = c.id
                   WHERE spp.status = true AND sp.status = true AND c.parent_id = ?`
		params = append(params, category)
	} else {
		query += ` WHERE spp.status = true AND sp.status = true`
	}

	// Add name filter for name_ka OR name_en if provided also dont make it to long
	if name != "" && len(name) < 255 {
		query += ` AND (LOWER(sp.name_ka) LIKE LOWER(?) OR LOWER(sp.name_en) LIKE LOWER(?) OR LOWER(sp.company) LIKE LOWER(?))`
		params = append(params, "%"+name+"%", "%"+name+"%", "%"+name+"%")
	}

	// Complete the query
	query += ` GROUP BY spp.product_id
               HAVING COUNT(*) > 1
              ) AS sub;`

	// Execute query
	err := r.db.Get(&count, query, params...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySQLStoreRepository) GetForCategorySlider(ctx context.Context) ([]models.CategorySlider, error) {
	var items []models.CategorySlider

	const (
		topCategories = `
		WITH top_categories AS (
			SELECT id, name
			FROM categories
			WHERE parent_id IS NULL
			ORDER BY RAND()
			LIMIT 3
		)`
		products = `,
		products AS (
			SELECT 
				sp.id,
				sp.name_ka AS name,
				sp.company,
				sp.image_url,
				sp.volume,
				sp.origin,
				MIN(spp.price) AS min_price,
				MAX(spp.price) AS max_price,
				c.parent_id,
				ROW_NUMBER() OVER (
					PARTITION BY c.parent_id
					ORDER BY (MAX(spp.price) - MIN(spp.price)) DESC
				) AS rn
			FROM store_products sp
			JOIN store_product_prices spp 
				ON sp.id = spp.product_id 
				AND spp.status = TRUE
			JOIN categories c 
				ON sp.category_id = c.id
			WHERE sp.status = TRUE
			GROUP BY 
				sp.id, 
				sp.name_ka, 
				sp.company, 
				sp.image_url, 
				sp.volume, 
				sp.origin, 
				c.parent_id
		)`
		mainQuery = `,
		SELECT 
			tc.name AS name,
			JSON_ARRAYAGG(
				JSON_OBJECT(
					'id', p.id,
					'name', p.name,
					'company', p.company,
					'image', p.image_url,
					'volume', p.volume,
					'origin', p.origin,
					'min_price', p.min_price,
					'max_price', p.max_price
				)
			) AS products
		FROM top_categories tc
		JOIN products p 
			ON p.parent_id = tc.id
		WHERE p.rn <= 6
		GROUP BY tc.name`
	)
	query := topCategories + products + mainQuery

	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("failed to get category slider products: %w", err)
	}

	return items, nil
}

func (r *MySQLStoreRepository) GetProductByBarCode(barCode string) (int64, error) {
	var productId int64

	const barCodeQuery = `SELECT product_id FROM store_product_bar_codes WHERE bar_code = ? LIMIT 1;`
	err := r.db.Get(&productId, barCodeQuery, barCode)

	if err != nil {
		return productId, err
	}

	return productId, nil
}

func (r *MySQLStoreRepository) CreateItem(data dto.Product) (int64, error) {
	metaJSON, err := json.Marshal(data.Meta)
	if err != nil {
		return 0, fmt.Errorf("failed to marshall: %w", err)
	}

	const insertQuery = `INSERT INTO store_products (name_ka, image_url, volume, meta) VALUES (?, ?, ?, ?);`
	res, err := r.db.Exec(insertQuery, data.Name, data.Image, &data.Volume, metaJSON)
	if err != nil {
		return 0, fmt.Errorf("failed to insert store product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	const insertBarCodeQuery = `INSERT INTO store_product_bar_codes (bar_code, product_id) VALUES (?, ?);`
	_, err = r.db.Exec(insertBarCodeQuery, data.BarCode, id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert bar code: %w", err)
	}

	return id, nil
}

func (r *MySQLStoreRepository) AddOrUpdatePrice(itemId int64, data dto.ProductPrice) error {
	const deactivateOldPrices = `UPDATE store_product_prices SET status = false WHERE status = true AND store_id = ? AND product_id = ? AND DATE(created_at) < CURDATE();`

	_, err := r.db.Exec(deactivateOldPrices, data.StoreId, itemId)
	if err != nil {
		return errors.New("failed to deactivate old prices")
	}

	const alreadyCreated = `SELECT id FROM store_product_prices  WHERE product_id = ? AND store_id = ? AND created_at > CURDATE();`

	var existingID string
	err = r.db.Get(&existingID, alreadyCreated, itemId, data.StoreId)
	if err == nil {
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	// Generate ULID
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	const insertQuery = `INSERT INTO store_product_prices (id, store_id, product_id, price, old_price) VALUES (?, ?, ?, ?, ?)`

	_, err = r.db.Exec(insertQuery, id.String(), data.StoreId, itemId, data.Price, data.OldPrice)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLStoreRepository) DisableOldPrices() error {
	const deactivateOldPrices = `UPDATE store_product_prices SET status = false WHERE status = true AND DATE(created_at) < DATE_SUB(CURDATE(), INTERVAL 3 DAY);`

	_, err := r.db.Exec(deactivateOldPrices)
	if err != nil {
		return errors.New("failed to deactivate old prices")
	}

	return nil
}
