package service

import (
	"errors"
	"strings"

	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
)

type dataCollectorConfigService struct {
	repo port.DataCollectorConfigRepoPort
}

func NewDataCollectorConfigService(repo port.DataCollectorConfigRepoPort) domain.DataCollectorConfigService {
	return &dataCollectorConfigService{repo: repo}
}

func (s dataCollectorConfigService) GetAllDataCollectorConfig() ([]port.VendorAccount, error) {
	dataCollectorConfig, err := s.repo.GetAllDataCollectorConfig()
	if err != nil {
		return []port.VendorAccount{}, err
	}
	return dataCollectorConfig, nil
}

func (s dataCollectorConfigService) GetDataCollectorConfigByVendorType(vendorType string) ([]port.VendorAccount, error) {
	dataCollectorConfig, err := s.repo.GetDataCollectorConfigByVendorType(vendorType)
	if err != nil {
		return []port.VendorAccount{}, err
	}
	return dataCollectorConfig, nil
}

func (s dataCollectorConfigService) GetOneDataCollectorConfig(id int) (port.VendorAccount, error) {
	dataCollectorConfig, err := s.repo.GetOneDataCollectorConfig(id)
	if err != nil {
		return port.VendorAccount{}, err
	}
	return dataCollectorConfig, nil
}

func (s dataCollectorConfigService) CreateDataCollectorConfig(dataCollectorConfig port.VendorAccount) error {
	dataCollectorConfig.VendorType = strings.ToLower(dataCollectorConfig.VendorType)
	switch dataCollectorConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		dataCollectorConfig.AppID = nil
		dataCollectorConfig.AppSecret = nil
	case constant.VENDOR_TYPE_KSTAR:
		dataCollectorConfig.AppID = nil
		dataCollectorConfig.AppSecret = nil
		dataCollectorConfig.Token = nil
	case constant.VENDOR_TYPE_INVT:
		dataCollectorConfig.Token = nil
	case constant.VENDOR_TYPE_HUAWEI:
		dataCollectorConfig.AppID = nil
		dataCollectorConfig.AppSecret = nil
		dataCollectorConfig.Token = nil
	default:
		err := errors.New("vendor type is incorrect")
		return err
	}

	err := s.repo.CreateDataCollectorConfig(dataCollectorConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s dataCollectorConfigService) UpdateDataCollectorConfig(id int, dataCollectorConfig port.VendorAccount) error {
	dataCollectorConfig.VendorType = strings.ToLower(dataCollectorConfig.VendorType)
	switch dataCollectorConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		dataCollectorConfig.AppID = nil
		dataCollectorConfig.AppSecret = nil
	case constant.VENDOR_TYPE_KSTAR:
		dataCollectorConfig.AppID = nil
		dataCollectorConfig.AppSecret = nil
		dataCollectorConfig.Token = nil
	case constant.VENDOR_TYPE_INVT:
		dataCollectorConfig.Token = nil
	case constant.VENDOR_TYPE_HUAWEI:
		dataCollectorConfig.AppID = nil
		dataCollectorConfig.AppSecret = nil
		dataCollectorConfig.Token = nil
	default:
		err := errors.New("vendor type is incorrect")
		return err
	}

	err := s.repo.UpdateDataCollectorConfig(id, dataCollectorConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s dataCollectorConfigService) DeleteDataCollectorConfig(id int) error {
	err := s.repo.DeleteDataCollectorConfig(id)
	if err != nil {
		return err
	}
	return nil
}
