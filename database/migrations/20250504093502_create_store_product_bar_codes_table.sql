-- +goose Up
-- +goose StatementBegin
CREATE TABLE store_product_bar_codes
(
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    bar_code   CHAR(30) NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_bar_code_product FOREIGN KEY (product_id) REFERENCES store_products (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE store_product_bar_codes;
-- +goose StatementEnd
