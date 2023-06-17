package service

import (
	"errors"
	"strings"

	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
)

type alarmConfigService struct {
	repo port.AlarmConfigRepoPort
}

func NewAlarmConfigService(repo port.AlarmConfigRepoPort) domain.AlarmConfigService {
	return &alarmConfigService{repo: repo}
}

func (s *alarmConfigService) GetAllAlarmConfig() ([]port.VendorAccount, error) {
	alarmConfig, err := s.repo.GetAllAlarmConfig()
	if err != nil {
		return []port.VendorAccount{}, err
	}
	return alarmConfig, nil
}

func (s *alarmConfigService) GetAlarmConfigByVendorType(vendorType string) ([]port.VendorAccount, error) {
	alarmConfig, err := s.repo.GetAlarmConfigByVendorType(vendorType)
	if err != nil {
		return []port.VendorAccount{}, err
	}
	return alarmConfig, nil
}

func (s *alarmConfigService) GetOneAlarmConfig(id int) (port.VendorAccount, error) {
	alarmConfig, err := s.repo.GetOneAlarmConfig(id)
	if err != nil {
		return port.VendorAccount{}, err
	}
	return alarmConfig, nil
}

func (s *alarmConfigService) CreateAlarmConfig(alarmConfig port.VendorAccount) error {
	alarmConfig.VendorType = strings.ToLower(alarmConfig.VendorType)
	switch alarmConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		alarmConfig.AppID = nil
		alarmConfig.AppSecret = nil
	case constant.VENDOR_TYPE_KSTAR:
		alarmConfig.AppID = nil
		alarmConfig.AppSecret = nil
		alarmConfig.Token = nil
	case constant.VENDOR_TYPE_INVT:
		alarmConfig.Token = nil
	case constant.VENDOR_TYPE_HUAWEI:
		alarmConfig.AppID = nil
		alarmConfig.AppSecret = nil
		alarmConfig.Token = nil
	default:
		err := errors.New("vendor type is incorrect")
		return err
	}

	err := s.repo.CreateAlarmConfig(alarmConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *alarmConfigService) UpdateAlarmConfig(id int, alarmConfig port.VendorAccount) error {
	alarmConfig.VendorType = strings.ToLower(alarmConfig.VendorType)
	switch alarmConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		alarmConfig.AppID = nil
		alarmConfig.AppSecret = nil
	case constant.VENDOR_TYPE_KSTAR:
		alarmConfig.AppID = nil
		alarmConfig.AppSecret = nil
		alarmConfig.Token = nil
	case constant.VENDOR_TYPE_INVT:
		alarmConfig.Token = nil
	case constant.VENDOR_TYPE_HUAWEI:
		alarmConfig.AppID = nil
		alarmConfig.AppSecret = nil
		alarmConfig.Token = nil
	default:
		err := errors.New("vendor type is incorrect")
		return err
	}

	err := s.repo.UpdateAlarmConfig(id, alarmConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *alarmConfigService) DeleteAlarmConfig(id int) error {
	err := s.repo.DeleteAlarmConfig(id)
	if err != nil {
		return err
	}
	return nil
}
