-- +goose Up
-- +goose StatementBegin
CREATE TABLE store_products
(
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    category_id BIGINT UNSIGNED NULL,
    name_ka     VARCHAR(255),
    name_en     VARCHAR(255),
    company     VARCHAR(255),
    image_url   VARCHAR(255),
    volume      VARCHAR(255) NULL,
    origin      VARCHAR(255) NULL,
    meta        JSON NULL,
    status      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_product_store FOREIGN KEY (category_id) REFERENCES categories (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE store_products;
-- +goose StatementEnd
