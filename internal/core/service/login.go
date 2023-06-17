package service

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"golang.org/x/crypto/bcrypt"
)

type loginService struct {
	repo               port.UserRepoPort
	secretKeyAuthToken string
}

func NewLoginService(repo port.UserRepoPort, secretKeyAuthToken string) domain.LoginService {
	return &loginService{repo: repo, secretKeyAuthToken: secretKeyAuthToken}
}

func (s loginService) Login(userCredential domain.Credential) (domain.AuthToken, int, error) {
	user, err := s.repo.GetUserByUsername(userCredential.Username)
	if err != nil {
		return domain.AuthToken{}, 0, err
	}

	err = s.ComparePassword(userCredential.Password, user.Password)
	if err != nil {
		return domain.AuthToken{}, user.ID, err
	}

	authToken, err := s.GenerateAuthToken(user)
	if err != nil {
		return domain.AuthToken{}, user.ID, err
	}

	return authToken, user.ID, nil
}

func (s loginService) GenerateAuthToken(user port.User) (domain.AuthToken, error) {
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = user.ID
	atClaims["expired_time"] = time.Now().UTC().Add(60 * time.Minute).Format(time.RFC3339Nano)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.secretKeyAuthToken))
	if err != nil {
		return domain.AuthToken{}, err
	}

	return domain.AuthToken{
		AccessToken: accessTokenString,
		ExpiredTime: atClaims["expired_time"].(string),
	}, nil
}

func (s loginService) GenerateHashPassword(rawPassword string) (string, error) {
	bytePassword := []byte(rawPassword)
	hashed, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s loginService) ComparePassword(password string, hashPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
