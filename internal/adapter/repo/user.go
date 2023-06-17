package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) port.UserRepoPort {
	return &userRepo{db: db}
}

func (r userRepo) GetUserByUsername(username string) (port.User, error) {
	queryString := `SELECT id, username, password, created_at, updated_at
					FROM user
					WHERE username = ?`
	row := r.db.QueryRow(queryString, username)
	if err := row.Err(); err != nil {
		return port.User{}, err
	}
	var result port.User
	err := row.Scan(
		&result.ID,
		&result.Username,
		&result.Password,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return port.User{}, err
	}

	return result, nil
}

func (r userRepo) GetKibanaUser() (port.KibanaUser, error) {
	queryString := `SELECT username, password
					FROM kibana_user
					WHERE username = ?
					LIMIT 1`
	row := r.db.QueryRow(queryString, "solar-viewer")
	if err := row.Err(); err != nil {
		return port.KibanaUser{}, err
	}
	var result port.KibanaUser
	err := row.Scan(
		&result.Username,
		&result.Password,
	)
	if err != nil {
		return port.KibanaUser{}, err
	}

	return result, nil
}

func (r userRepo) CreateUser(user port.User) error {
	queryString := `INSERT INTO user (username, password)
					VALUES (?, ?)`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		user.Username,
		user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}
