package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type userService struct {
	repo port.UserRepoPort
}

func NewUserService(repo port.UserRepoPort) domain.UserService {
	return &userService{repo: repo}
}

func (s userService) GetUserByUsername(username string) (port.User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return port.User{}, err
	}

	return user, nil
}

func (s userService) GetKibanaUser() (port.KibanaUser, error) {
	user, err := s.repo.GetKibanaUser()
	if err != nil {
		return port.KibanaUser{}, err
	}

	return user, nil
}
