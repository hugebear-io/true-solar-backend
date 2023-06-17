package port

import "time"

type UserRepoPort interface {
	GetUserByUsername(username string) (User, error)
	GetKibanaUser() (KibanaUser, error)
	CreateUser(user User) error
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KibanaUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
