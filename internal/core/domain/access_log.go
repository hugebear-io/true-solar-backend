package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type AccessLogService interface {
	GetAccessLog(limit int, offset int) ([]port.AccessLog, int, error)
	CreateAccessLog(accessLog port.AccessLog) error
}
