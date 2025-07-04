-- +goose Up
-- +goose StatementBegin
CREATE TABLE currency_providers
(
    id       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name     VARCHAR(255),
    logo_url VARCHAR(255),
    PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE currency_providers;
-- +goose StatementEnd
