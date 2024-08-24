package domain

type User struct {
	Email           string
	Password        string
	ConfirmPassword string
	CreateTime      int64
}
