package repository

import (
	"context"
	"fmt"

	"github.com/patrickn2/api-challenge/infra/db"
	"github.com/patrickn2/api-challenge/interfaces"
	"github.com/patrickn2/api-challenge/schema"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *db.Database
}

func NewUserRepository(db *db.Database) interfaces.UserRepository {
	return &UserRepository{db}
}

func (p *UserRepository) InsertUsers(ctx context.Context, users []*schema.User) (int, error) {
	db := p.db.Conn.WithContext(ctx)
	result := db.CreateInBatches(users, 200)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (p *UserRepository) GetClerks(ctx context.Context, params *schema.GetClerksParams) ([]*schema.User, error) {
	// Limit will be (limit + 1) to check if there is a next or previous page
	db := p.db.Conn.WithContext(ctx)
	if params.StartingAfter != nil || params.EndingBefore != nil {
		return getStartEndingUsers(db, params)
	}

	var users []*schema.User
	tx := db.Model(&schema.User{})
	tx = filterEmail(tx, params.Email)
	tx.Order("created_at DESC").Limit(int(*params.Limit) + 1).Find(&users)
	return users, nil
}

func filterEmail(tx *gorm.DB, email *string) *gorm.DB {
	if email != nil {
		return tx.Where("email ~ ?", *email)
	}
	return tx
}

func getStartEndingUsers(db *gorm.DB, params *schema.GetClerksParams) ([]*schema.User, error) {
	var users []*schema.User
	var v struct {
		createdAtSignal string
		createdAtOrder  string
		userId          uint
		reverse         bool
	}
	if params.StartingAfter != nil {
		v.userId = *params.StartingAfter
		v.createdAtSignal = "<"
		v.createdAtOrder = "DESC"
		v.reverse = false
	} else {
		v.userId = *params.EndingBefore
		v.createdAtSignal = ">"
		v.createdAtOrder = "ASC"
		v.reverse = true
	}
	subQueryId := db.Select("created_at").Model(&schema.User{}).Where("id = ?", v.userId)
	query := db.Model(&schema.User{}).Where(fmt.Sprintf("created_at %s (?)", v.createdAtSignal), subQueryId).Limit(int(*params.Limit) + 1).Order(fmt.Sprintf("created_at %s", v.createdAtOrder))
	if params.Email != nil {
		query = query.Where("email ~ ?", *params.Email)
	}
	if v.reverse {
		query = db.Table("(?) as u", query).Order("u.created_at DESC")
	}
	query.Find(&users)
	return users, nil
}
