package repository

import (
	"awesome-bluebook/domain"
	"awesome-bluebook/repository/dao"
	"context"
	"errors"
)

var _ UserRepository = (*BasicUserRepository)(nil)
var UserDuplicateErr = dao.UserDuplicateErr

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type BasicUserRepository struct {
	dao dao.UserDAO
}

func NewBasicUserRepository(dao dao.UserDAO) UserRepository {
	return &BasicUserRepository{
		dao: dao,
	}
}

func (repo *BasicUserRepository) Create(ctx context.Context, user domain.User) error {
	err := repo.dao.Insert(ctx, domainToEntity(user))
	switch {
	case errors.Is(err, dao.UserDuplicateErr):
		return UserDuplicateErr
	default:
		return err
	}
}

func domainToEntity(user domain.User) dao.User {
	return dao.User{
		Email:    user.Email,
		Password: user.Password,
	}
}
