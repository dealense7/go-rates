package repositories

import (
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type MySQLGasRepository struct {
	db *sqlx.DB
}

func NewMySQLGasRepository(db *sqlx.DB) interfaces.GasRepository {
	return &MySQLGasRepository{db: db}
}

func (r *MySQLGasRepository) FindAll() ([]models.Gas, error) {
	var items []models.Gas
	query := `select gr.id, gr.name, gr.tag, gr.price, gp.name as provider_name, gp.logo_url as provider_logo
				from gas_rates as gr
						 join gas_providers gp on gp.id = gr.provider_id
				where gr.status is true;`

	err := r.db.Select(&items, query)
	return items, err
}
