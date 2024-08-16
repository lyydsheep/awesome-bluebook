package repository

import "github.com/lyydsheep/awesome-bluebook/internal/domain"

type BasicUserRepository interface {
	Create(u domain.User) error
}

type Repository struct {
}

func (r *Repository) Create(u domain.User) error {
	//TODO implement me
	panic("implement me")
}
