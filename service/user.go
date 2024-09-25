package service

import (
	"awesome-bluebook/domain"
	"context"
)

var _ UserService = (*BasicUserService)(nil)

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
}

type BasicUserService struct {
	repo UserRepository
}

func (svc *BasicUserService) Signup(ctx context.Context, user domain.User) error {
	return svc.repo.Create(ctx, user)
}
