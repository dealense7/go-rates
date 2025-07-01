package repositories

import (
	"github.com/dealense7/go-rate-app/internal/interfaces"
	"github.com/dealense7/go-rate-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type MySQLCategoryRepository struct {
	db *sqlx.DB
}

func NewMySQLCategoryRepository(db *sqlx.DB) interfaces.CategoryRepository {
	return &MySQLCategoryRepository{db: db}
}

func (r *MySQLCategoryRepository) GetItems() ([]models.Category, error) {
	var parents []models.Category
	var children []models.Category

	// Select Parents
	query := `SELECT * FROM categories WHERE parent_id IS NULL`
	if err := r.db.Select(&parents, query); err != nil {
		return nil, err
	}

	// Select Children
	query = `SELECT * FROM categories WHERE parent_id IS NOT NULL`
	if err := r.db.Select(&children, query); err != nil {
		return nil, err
	}

	// Map children by parent_id
	childMap := make(map[int64][]models.Category)
	for _, child := range children {
		if child.ParentID != nil {
			childMap[*child.ParentID] = append(childMap[*child.ParentID], child)
		}
	}

	// Assign children to their respective parents
	for _, parent := range parents {
		parent.Children = childMap[parent.ID]
	}

	return parents, nil
}
