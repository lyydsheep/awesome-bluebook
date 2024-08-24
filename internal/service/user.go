package service

import (
	"context"
	"github.com/lyydsheep/awesome-bluebook/internal/domain"
	"github.com/lyydsheep/awesome-bluebook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type BasicUserService interface {
	SignUp(ctx context.Context, u domain.User) error
}

type UserService struct {
	repo repository.BasicUserRepository
}

var ErrDuplicateUser = repository.ErrDuplicateUser

func NewUserService(repo repository.BasicUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 对密码进行加密
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	return s.repo.Create(ctx, u)
}
