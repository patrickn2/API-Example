package schema

import (
	"time"
)

type Picture struct {
	Thumbnail string `json:"thumbnail"`
}

type Registered struct {
	Date time.Time `json:"date"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type Result struct {
	Name       Name       `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Picture    Picture    `json:"picture"`
	Registered Registered `json:"registered"`
}

type RandomUser struct {
	Results []Result `json:"results"`
}

type User struct {
	Id        uint      `json:"id" gorm:"->;primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}
