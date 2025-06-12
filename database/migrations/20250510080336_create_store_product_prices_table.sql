-- +goose Up
-- +goose StatementBegin
CREATE TABLE store_product_prices
(
    id         CHAR(26) PRIMARY KEY,
    store_id   BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    price      BIGINT NOT NULL,
    old_price  BIGINT NOT NULL,
    status     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_store FOREIGN KEY (store_id) REFERENCES store_providers (id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES store_products (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_product_status_created ON store_product_prices (product_id, status, created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_store_product_product_status ON store_product_prices (product_id, status);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_store_product_prices_product_status ON store_product_prices (product_id, price, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE store_product_prices;
-- +goose StatementEnd
