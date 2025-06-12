package seeders

import (
	"fmt"
	"github.com/dealense7/go-rate-app/internal/enum"
	"github.com/jmoiron/sqlx"
)

func seedGasData(db *sqlx.DB) {
	items := []enum.GasProvider{
		enum.SOCAR,
		enum.WISSOL,
		enum.PORTAL,
		enum.GULF,
		enum.ROMPETROL,
		enum.LUKOILI,
		enum.CONNECT,
	}

	for _, p := range items {
		tx := db.MustBegin()

		const query = `INSERT INTO gas_providers (id, name, slug, logo_url)
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
