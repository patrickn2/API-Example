package model

import (
	schema "github.com/patrickn2/Clerk-Challenge/schemas"
	"gorm.io/gorm"
)

func InsertNewUsers(db *gorm.DB, users *[]schema.User) (int, error) {
	result := db.CreateInBatches(users, 200)
	if result.Error != nil {
		return int(result.RowsAffected), result.Error
	}
	return int(result.RowsAffected), nil
}
