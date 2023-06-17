package service

import (
	"errors"

	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type performanceAlarmConfigService struct {
	repo port.PerformanceAlarmConfigRepoPort
}

func NewPerformanceAlarmConfigService(repo port.PerformanceAlarmConfigRepoPort) domain.PerformanceAlarmConfigService {
	return &performanceAlarmConfigService{repo: repo}
}

func (s performanceAlarmConfigService) GetPerformanceAlarmConfig() (domain.UpdatePerformanceAlarmRequestBody, error) {
	performanceAlarmConfig, err := s.repo.GetPerformanceAlarmConfig()
	if err != nil {
		return domain.UpdatePerformanceAlarmRequestBody{}, err
	}

	if len(performanceAlarmConfig) != 2 {
		err := errors.New("performance alarm configuration must have two items")
		return domain.UpdatePerformanceAlarmRequestBody{}, err
	}

	performanceAlarm := domain.UpdatePerformanceAlarmRequestBody{
		Alarm1ID:         performanceAlarmConfig[0].ID,
		Alarm1Name:       performanceAlarmConfig[0].Name,
		Alarm1Interval:   performanceAlarmConfig[0].Interval,
		Alarm1HitDay:     performanceAlarmConfig[0].HitDay,
		Alarm1Percentage: performanceAlarmConfig[0].Percentage,
		Alarm1Duration:   performanceAlarmConfig[0].Duration,
		Alarm1CreatedAt:  performanceAlarmConfig[0].CreatedAt,
		Alarm1UpdatedAt:  performanceAlarmConfig[0].UpdatedAt,
		Alarm2ID:         performanceAlarmConfig[1].ID,
		Alarm2Name:       performanceAlarmConfig[1].Name,
		Alarm2Interval:   performanceAlarmConfig[1].Interval,
		Alarm2Percentage: performanceAlarmConfig[1].Percentage,
		Alarm2Duration:   performanceAlarmConfig[1].Duration,
		Alarm2CreatedAt:  performanceAlarmConfig[1].CreatedAt,
		Alarm2UpdatedAt:  performanceAlarmConfig[1].UpdatedAt,
	}

	return performanceAlarm, nil
}

func (s performanceAlarmConfigService) UpdatePerformanceAlarmConfig(performanceAlarmConfig domain.UpdatePerformanceAlarmRequestBody) error {
	// existingConfig, err := s.repo.GetPerformanceAlarmConfig()
	// if err != nil {
	// 	return err
	// }

	// if len(existingConfig) != 2 {
	// 	err := errors.New("performance alarm configuration must have two items")
	// 	return err
	// }

	// toBeUpdatedPerformanceLow := domain.PerformanceAlarmConfig{
	// 	Interval:   performanceAlarmConfig.Alarm1Interval,
	// 	HitDay:     &performanceAlarmConfig.Alarm1HitDay,
	// 	Percentage: performanceAlarmConfig.Alarm1Percentage,
	// 	Duration:   &performanceAlarmConfig.Alarm1Duration,
	// }
	// err = s.repo.UpdatePerformanceAlarmConfig(1, toBeUpdatedPerformanceLow)
	// if err != nil {
	// 	return err
	// }

	// if toBeUpdatedPerformanceLow.Interval != existingConfig[0].Interval ||
	// 	toBeUpdatedPerformanceLow.Percentage != existingConfig[0].Percentage ||
	// 	(pointy.IntValue(toBeUpdatedPerformanceLow.HitDay, -1) != pointy.IntValue(existingConfig[0].HitDay, -1)) ||
	// 	(pointy.IntValue(toBeUpdatedPerformanceLow.Duration, -1) != pointy.IntValue(existingConfig[0].Duration, -1)) {

	// 	updatedConfig, err := s.GetPerformanceAlarmConfig()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	go func() {
	// 		s.performanceLowConfigCh <- domain.PerformanceAlarmConfig{
	// 			ID:         updatedConfig.Alarm1ID,
	// 			Name:       updatedConfig.Alarm1Name,
	// 			Interval:   updatedConfig.Alarm1Interval,
	// 			HitDay:     &updatedConfig.Alarm1HitDay,
	// 			Percentage: updatedConfig.Alarm1Percentage,
	// 			Duration:   &updatedConfig.Alarm1Duration,
	// 			CreatedAt:  updatedConfig.Alarm1CreatedAt,
	// 			UpdatedAt:  updatedConfig.Alarm1UpdatedAt,
	// 		}
	// 	}()
	// }

	// toBeUpdatedSumPerformanceLow := domain.PerformanceAlarmConfig{
	// 	Interval:   performanceAlarmConfig.Alarm2Interval,
	// 	Percentage: performanceAlarmConfig.Alarm2Percentage,
	// 	Duration:   &performanceAlarmConfig.Alarm2Duration,
	// }
	// err = s.repo.UpdatePerformanceAlarmConfig(2, toBeUpdatedSumPerformanceLow)
	// if err != nil {
	// 	return err
	// }

	// if toBeUpdatedSumPerformanceLow.Interval != existingConfig[1].Interval ||
	// 	toBeUpdatedSumPerformanceLow.Percentage != existingConfig[1].Percentage ||
	// 	(pointy.IntValue(toBeUpdatedSumPerformanceLow.Duration, -1) != pointy.IntValue(existingConfig[1].Duration, -1)) {

	// 	updatedConfig, err := s.GetPerformanceAlarmConfig()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	go func() {
	// 		s.sumPerformanceLowConfigCh <- domain.PerformanceAlarmConfig{
	// 			ID:         updatedConfig.Alarm2ID,
	// 			Name:       updatedConfig.Alarm2Name,
	// 			Interval:   updatedConfig.Alarm2Interval,
	// 			Percentage: updatedConfig.Alarm2Percentage,
	// 			Duration:   &updatedConfig.Alarm2Duration,
	// 			CreatedAt:  updatedConfig.Alarm2CreatedAt,
	// 			UpdatedAt:  updatedConfig.Alarm2UpdatedAt,
	// 		}
	// 	}()
	// }

	// return nil
	return nil
}
