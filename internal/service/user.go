package service

import "github.com/lyydsheep/awesome-bluebook/internal/domain"

type BasicUserService interface {
	SignUp(u domain.User) error
}

type UserService struct {
}

func (s *UserService) SignUp(u domain.User) error {
	//TODO implement me
	panic("implement me")
}
