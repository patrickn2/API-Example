package schema

import (
	"time"
)

type User struct {
	Id        uint      `json:"id" gorm:"->;primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}

type GetClerksParams struct {
	Limit         *uint   `form:"limit"`
	StartingAfter *uint   `form:"starting_after"`
	EndingBefore  *uint   `form:"ending_before"`
	Email         *string `form:"email"`
}
