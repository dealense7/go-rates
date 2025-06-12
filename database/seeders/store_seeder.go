package seeders

import (
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/jmoiron/sqlx"
)

func seedStoreData(db *sqlx.DB) {
	items := []enum.StoreProvider{
		enum.GOODWILL,
		enum.CARREFOUR,
		enum.ORINABIJI,
		enum.AGROHUB,
		enum.MAGNITI,
	}

	for _, p := range items {
		tx := db.MustBegin()

		const query = `INSERT INTO store_providers (id, name, slug, logo_url)
						VALUES (?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE name     = VALUES(name),
												slug     = VALUES(slug),
												logo_url = VALUES(logo_url)`

		tx.MustExec(query, int(p), p.String(), p.Slug(), p.Logo())
		err := tx.Commit()
		if err != nil {
			fmt.Println(err.Error())
		}

	}
}
