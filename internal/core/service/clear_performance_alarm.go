package service

import (
	"fmt"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type clearPerformanceAlarmService struct {
	elastic port.ElasticSearchRepoPort
}

func NewClearPerformanceAlarmService(elastic port.ElasticSearchRepoPort) *clearPerformanceAlarmService {
	return &clearPerformanceAlarmService{elastic: elastic}
}

func (s *clearPerformanceAlarmService) Run() {
	result, _ := s.elastic.QueryPerformanceLow(3, 60, 24, 60/100)
	fmt.Printf("%#v\n", result)
}
