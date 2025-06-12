package repositories

import (
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

	err := r.db.Get(&item, query, id)
	if err != nil {
		return item, err
	}
	fmt.Printf("%#v\n", item)

	return item, nil
}

func (r *MySQLStoreRepository) GetForSlider() ([]models.ProductItem, error) {
	var items []models.ProductItem

	query := `WITH interesting_products AS (SELECT product_id
											  FROM (SELECT product_id
													FROM store_product_prices
													JOIN store_products ON store_products.id = product_id and store_products.status = true
													WHERE store_product_prices.status = true
													GROUP BY product_id
													HAVING COUNT(product_id) > 1
													ORDER BY (MAX(price) - MIN(price)) DESC
													limit 50
													) AS top_products
 											  ORDER BY RAND()
 											  limit 6
											  )
				SELECT sp.id,
					   sp.name_ka AS name,
					   sp.company,
					   sp.image_url,
					   sp.volume,
					   sp.origin,
					   min(spp.price) as min_price,
					   max(spp.price) as max_price
				FROM store_products AS sp
						 JOIN interesting_products ip ON ip.product_id = sp.id
						 JOIN golang.store_product_prices spp
							  ON spp.product_id = sp.id AND spp.status = true
				GROUP BY sp.id LIMIT 18;`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *MySQLStoreRepository) GetItemsList(offset int) ([]models.ProductItem, error) {
	var items []models.ProductItem

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
						  ON spp.product_id = sp.id AND spp.status = true
			WHERE sp.status = true
			GROUP BY sp.id, sp.name_ka, sp.company, sp.image_url, sp.volume, sp.origin
			having count(spp.product_id) > 1
			order by (max_price - min_price) desc LIMIT 30 OFFSET ?;`

	err := r.db.Select(&items, query, offset)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *MySQLStoreRepository) GetItemsCount() (int, error) {
	var count int

	query := `SELECT COUNT(*)
				FROM (
						 SELECT spp.product_id
						 FROM store_product_prices AS spp
								  JOIN store_products sp ON sp.id = spp.product_id
						 WHERE spp.status = true AND sp.status = true
						 GROUP BY spp.product_id
						 HAVING COUNT(*) > 1
					 ) AS sub;
				;`

	err := r.db.Get(&count, query)

	return count, err
}

func (r *MySQLStoreRepository) GetForCategorySlider() ([]models.CategorySlider, error) {
	var items []models.CategorySlider

	query := `WITH top_categories AS (
				SELECT id, name
				FROM categories
				WHERE parent_id IS NULL
				ORDER BY RAND()
				LIMIT 3
			),
				 products_with_prices AS (
					 SELECT
						 sp.id,
						 sp.name_ka AS name,
						 sp.company,
						 sp.image_url,
						 sp.volume,
						 sp.origin,
						 spp.price,
						 sp.category_id,
						 c.parent_id
					 FROM store_products sp
							  JOIN store_product_prices spp
								   ON sp.id = spp.product_id AND spp.status = TRUE
							  JOIN categories c
								   ON sp.category_id = c.id
					 WHERE sp.status = TRUE
				 ),
				 min_max_prices AS (
					 SELECT
						 id,
						 MIN(price) AS min_price,
						 MAX(price) AS max_price
					 FROM products_with_prices
					 GROUP BY id
				 ),
				 filtered_products AS (
					 SELECT
						 p.id,
						 p.name,
						 p.company,
						 p.image_url,
						 p.volume,
						 p.origin,
						 m.min_price,
						 m.max_price,
						 tc.id AS top_category_id,
						 tc.name AS top_category_name
					 FROM products_with_prices p
							  JOIN min_max_prices m ON p.id = m.id
							  JOIN top_categories tc ON tc.id = p.parent_id
					 GROUP BY p.id, tc.id, tc.name
				 ),
				 ranked_products AS (
					 SELECT *,
							ROW_NUMBER() OVER (
								PARTITION BY top_category_id
								ORDER BY (max_price - min_price) DESC
								) AS rn
					 FROM filtered_products
				 )
			SELECT
				top_category_name as name,
				JSON_ARRAYAGG(
						JSON_OBJECT(
								'id', id,
								'name', name,
								'company', company,
								'image', image_url,
								'volume', volume,
								'origin', origin,
								'min_price', min_price,
								'max_price', max_price
						)
				) AS products
			FROM ranked_products
			WHERE rn <= 6
			GROUP BY top_category_name;`

	err := r.db.Select(&items, query)
	if err != nil {
		return nil, err
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
