package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type LoginService interface {
	Login(userCredential Credential) (AuthToken, int, error)
	GenerateAuthToken(user port.User) (AuthToken, error)
	GenerateHashPassword(rawPassword string) (string, error)
	ComparePassword(password string, hashPassword string) error
}

type Credential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiredTime string `json:"expired_time"`
}
