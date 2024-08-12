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
	emailQuery := "email LIKE ?"
	l, limitErr := strconv.Atoi(limit)
	if limitErr != nil {
		l = 10
	}
	if l > 100 {
		l = 100
	}
	if l < 1 {
		l = 10
	}

	email = strings.ToLower(strings.ReplaceAll(email, "%", ""))
	var users []schema.User
	tx := db.Select("id, name, email, phone, picture, created_at")
	if startingAfter != "" {
		sa, _ := strconv.Atoi(startingAfter)
		subQueryId := db.Select("created_at").Model(&schema.User{}).Where("id = ?", sa)
		tx = tx.Model(&schema.User{}).Where("created_at < (?)", subQueryId).Limit(l).Order("created_at DESC")
		if email != "" {
			tx = tx.Where(emailQuery, email+"%")
		}
		tx.Find(&users)
		return &users, nil
	}

	if endingBefore != "" {
		eb, _ := strconv.Atoi(endingBefore)
		subQueryId := db.Select("created_at").Model(&schema.User{}).Where("id = ?", eb)
		subQueryBefore := db.Select("id, name, email, phone, picture, created_at").Model(&schema.User{}).Where("created_at > (?)", subQueryId).Order("created_at ASC").Limit(l)
		if email != "" {
			subQueryBefore = subQueryBefore.Where(emailQuery, email+"%")
		}
		tx.Table("(?) as u", subQueryBefore).Order("u.created_at DESC").Find(&users)
		return &users, nil
	}

	if email != "" {
		tx = tx.Where(emailQuery, email+"%")
	}
	tx.Order("created_at DESC").Limit(l).Find(&users)

	return &users, nil
}
