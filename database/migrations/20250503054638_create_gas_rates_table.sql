-- +goose Up
-- +goose StatementBegin
CREATE TABLE gas_rates (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    provider_id BIGINT UNSIGNED NOT NULL,
    name        VARCHAR(255)     NOT NULL,
    tag         VARCHAR(255)     NOT NULL,
    price       BIGINT           NOT NULL,
    date        TIMESTAMP        NOT NULL,
    status      BOOLEAN          NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP        NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_gas_rates_provider_name_date (provider_id, name, date),
    CONSTRAINT fk_gas_rates_provider FOREIGN KEY (provider_id)
        REFERENCES gas_providers (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gas_rates;
-- +goose StatementEnd
