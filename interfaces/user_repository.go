package interfaces

import (
	"context"

	"github.com/patrickn2/api-challenge/schema"
)

type UserRepository interface {
	InsertUsers(context.Context, []*schema.User) (int, error)
	GetClerks(context.Context, *schema.GetClerksParams) ([]*schema.User, error)
}
