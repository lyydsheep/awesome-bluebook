package domain

import "time"

type User struct {
	Email           string
	Password        string
	ConfirmPassword string
	CreateTime      time.Time
}
