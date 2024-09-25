package repository

import (
	"awesome-bluebook/domain"
	"awesome-bluebook/repository/dao"
	"context"
)

var _ UserRepository = (*BasicUserRepository)(nil)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type BasicUserRepository struct {
	dao UserDAO
}

func NewBasicUserRepository(dao UserDAO) {
	return &BasicUserRepository{
		dao: dao,
	}
}

func (repo *BasicUserRepository) Create(ctx context.Context, user domain.User) error {
}

func domainToEntity(user domain.User) dao.User {
	return dao.User{
		Email:    user.Email,
		Password: user.Password,
	}
}
