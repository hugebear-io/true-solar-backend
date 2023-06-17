package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type accessLogService struct {
	repo port.AccessLogRepoPort
}

func NewAccessLogService(repo port.AccessLogRepoPort) domain.AccessLogService {
	return &accessLogService{
		repo: repo,
	}
}

func (s accessLogService) GetAccessLog(limit int, offset int) ([]port.AccessLog, int, error) {
	total, err := s.repo.TotalAccessLog()
	if err != nil {
		return []port.AccessLog{}, 0, err
	}

	accessLog, err := s.repo.GetAccessLog(limit, offset)
	if err != nil {
		return []port.AccessLog{}, total, err
	}

	return accessLog, total, nil
}

func (s accessLogService) CreateAccessLog(accessLog port.AccessLog) error {
	err := s.repo.CreateAccessLog(accessLog)
	if err != nil {
		return err
	}
	return nil
}
