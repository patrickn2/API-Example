package model

import (
	"strconv"
	"strings"

	"github.com/patrickn2/Clerk-Challenge/schema"
	"gorm.io/gorm"
)

func InsertNewUsers(db *gorm.DB, users *[]schema.User) (int, error) {
	result := db.CreateInBatches(users, 200)
	if result.Error != nil {
		return int(result.RowsAffected), result.Error
	}
	return int(result.RowsAffected), nil
}

func GetUsers(db *gorm.DB, limit string, startingAfter string, endingBefore string, email string) (*[]schema.User, error) {
	if limit == "" {
		limit = "10"
	}
	l, _ := strconv.Atoi(limit)
	var users []schema.User
	tx := db.Select("id, name, email, phone, picture, created_at")
	if startingAfter != "" {
		sa, _ := strconv.Atoi(startingAfter)
		sub := db.Select("created_at").Table("random_project.users").Where("id = ?", sa)
		tx = tx.Table("random_project.users").Where("created_at < (?)", sub)
	} else if endingBefore != "" {
		eb, _ := strconv.Atoi(endingBefore)
		sub := db.Select("created_at").Table("random_project.users").Where("id = ?", eb)
		sub2 := db.Select("id, name, email, phone, picture, created_at").Table("random_project.users").Where("created_at > (?)", sub).Order("created_at ASC").Limit(l)
		tx = tx.Table("(?)", sub2).Order("created_at DESC")
	}
	if email != "" {

		tx = tx.Where("email LIKE '?%'", strings.ToLower(email))
	}

	tx.Order("created_at DESC").Limit(l).Find(&users)
	return &users, nil
}
