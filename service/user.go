package service

import (
	"awesome-bluebook/domain"
	"awesome-bluebook/repository"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var _ UserService = (*BasicUserService)(nil)
var UserDuplicateErr = repository.UserDuplicateErr

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
}

type BasicUserService struct {
	repo repository.UserRepository
	l    *zap.Logger
}

func NewBasicUserService(repo repository.UserRepository, l *zap.Logger) UserService {
	return &BasicUserService{
		repo: repo,
		l:    l,
	}
}

func (svc *BasicUserService) Signup(ctx context.Context, user domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		svc.l.Error("bcrypt加密错误", zap.Error(err))
		return err
	}
	user.Password = string(encrypted)
	err = svc.repo.Create(ctx, user)
	switch {
	case errors.Is(err, repository.UserDuplicateErr):
		return UserDuplicateErr
	default:
		return err
	}
}
