-- +goose Up
-- +goose StatementBegin
CREATE TABLE store_providers
(
    id       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name     VARCHAR(255),
    logo_url VARCHAR(255),
    slug     VARCHAR(255),
    PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE store_providers;
-- +goose StatementEnd
