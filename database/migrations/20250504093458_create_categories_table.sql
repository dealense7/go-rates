-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories
(
    id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    parent_id BIGINT UNSIGNED NULL,
    name      VARCHAR(255),
    PRIMARY KEY (id),
    CONSTRAINT fk_parent_id FOREIGN KEY (parent_id) REFERENCES categories (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories;
-- +goose StatementEnd
