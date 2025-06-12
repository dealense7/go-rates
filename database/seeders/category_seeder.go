package seeders

import (
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/jmoiron/sqlx"
)

func seedCategoryData(db *sqlx.DB) {
	items := []enum.StoreCategoryProvider{
		enum.Grocery,
		enum.Drinks,
		enum.Dairy,
		enum.Sweet,
	}

	for _, p := range items {
		tx := db.MustBegin()

		const query = `INSERT INTO categories (id, name)
						VALUES (?, ?)
						ON DUPLICATE KEY UPDATE name = VALUES(name)`

		tx.MustExec(query, int(p), p.String())
		err := tx.Commit()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
