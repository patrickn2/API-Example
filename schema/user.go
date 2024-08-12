package schema

import (
	"time"
)

type RandomUser struct {
	Results []struct {
		Name struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Picture struct {
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		Registered struct {
			Date time.Time `json:"date"`
		} `json:"registered"`
	} `json:"results"`
}

type User struct {
	Id        uint      `json:"id" gorm:"->;primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}
