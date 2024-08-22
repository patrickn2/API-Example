package interfaces

import "github.com/patrickn2/api-challenge/schema"

type GetClerksParams struct {
	Limit         *uint   `form:"limit"`
	StartingAfter *uint   `form:"starting_after"`
	EndingBefore  *uint   `form:"ending_before"`
	Email         *string `form:"email"`
}

type UserRepository interface {
	InsertUsers([]*schema.User) (int, error)
	GetClerks(*GetClerksParams) ([]*schema.User, error)
}
