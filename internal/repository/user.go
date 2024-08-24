package repository

import (
	"context"
	"database/sql"
	"github.com/lyydsheep/awesome-bluebook/internal/domain"
	"github.com/lyydsheep/awesome-bluebook/internal/repository/dao"
)

type BasicUserRepository interface {
	Create(ctx context.Context, u domain.User) error
	EntityToDomain(u dao.User) domain.User
	DomainToEntity(u domain.User) dao.User
}

type Repository struct {
	dao dao.BasicUserDAO
}

var ErrDuplicateUser = dao.ErrDuplicateUser

func (r *Repository) EntityToDomain(u dao.User) domain.User {
	return domain.User{
		Email:      u.Email.String,
		Password:   u.Password,
		CreateTime: u.CreateTime,
	}
}

func (r *Repository) DomainToEntity(u domain.User) dao.User {
	return dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email == "",
		},
		Password:   u.Password,
		CreateTime: u.CreateTime,
	}
}

func NewRepository(dao dao.BasicUserDAO) *Repository {
	return &Repository{
		dao: dao,
	}
}

func (r *Repository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.DomainToEntity(u))
}
