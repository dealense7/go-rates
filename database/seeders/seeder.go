package seeders

import "github.com/dealense7/go-rate-app/internal/utils"

func init() {
	db := utils.NewDB()

	seedGasData(db)
	seedStoreData(db)
	seedCategoryData(db)
}
