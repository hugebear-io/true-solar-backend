package port

import "time"

type AccessLogRepoPort interface {
	GetAccessLog(limit int, offset int) ([]AccessLog, error)
	CreateAccessLog(accessLog AccessLog) error
	TotalAccessLog() (int, error)
}

type AccessLog struct {
	ID         int       `json:"id"`
	Message    string    `json:"message"`
	ByUserID   int       `json:"by_user_id"`
	ByUsername string    `json:"by_username"`
	CreatedAt  time.Time `json:"created_at"`
}
