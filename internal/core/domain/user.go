package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type UserService interface {
	GetUserByUsername(username string) (port.User, error)
	GetKibanaUser() (port.KibanaUser, error)
}
