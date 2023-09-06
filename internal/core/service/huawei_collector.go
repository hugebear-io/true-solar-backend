package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/adapter/collector"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/huawei"
)

type huaweiCollectorService struct {
	dataCollectorConfig domain.DataCollectorConfigService
	siteRegionMapping   domain.SiteRegionMappingService
	elastic             port.ElasticSearchRepoPort
	logger              logger.Logger
}

func NewHuaweiCollectorService(
	dataCollectorConfig domain.DataCollectorConfigService,
	siteRegionMapping domain.SiteRegionMappingService,
	elastic port.ElasticSearchRepoPort,
	logger logger.Logger,
) *huaweiCollectorService {
	return &huaweiCollectorService{
		dataCollectorConfig: dataCollectorConfig,
		siteRegionMapping:   siteRegionMapping,
		elastic:             elastic,
		logger:              logger,
	}
}

func (s huaweiCollectorService) Run() error {
	documentCh := make(chan interface{})
	errorCh := make(chan error)
	doneCh := make(chan bool)

	documentBatches := make([]interface{}, 0)
	siteDocumentBatches := make([]port.SiteItem, 0)

	siteRegions, err := s.siteRegionMapping.GetAllSiteRegionMapping()
	if err != nil {
		s.logger.Error(err)
		return err
	}

	usernames := make([]string, 0)
	password := ""
	configs, err := s.dataCollectorConfig.GetDataCollectorConfigByVendorType(huawei.BRAND)
	if err != nil {

		s.logger.Error(err)
		return err
	}

	for _, config := range configs {
		usernames = append(usernames, config.Username)
		password = config.Password
	}

	collector := collector.NewHuaweiCollector(
		s.dataCollectorConfig,
		s.siteRegionMapping,
		s.logger,
		documentCh,
		errorCh,
		doneCh,
		usernames,
		password,
		siteRegions,
	)

	go collector.Run()

DONE:
	for {
		select {
		case <-doneCh:
			break DONE
		case err := <-errorCh:
			s.logger.Error(err)
			return err
		case doc := <-documentCh:
			documentBatches = append(documentBatches, doc)
			if plantItem, ok := doc.(port.PlantItem); ok {
				tmp := port.SiteItem{
					Timestamp:   &plantItem.Timestamp,
					VendorType:  plantItem.VendorType,
					Area:        plantItem.Area,
					SiteID:      plantItem.SiteID,
					NodeType:    plantItem.NodeType,
					Name:        plantItem.Name,
					Location:    plantItem.Location,
					PlantStatus: plantItem.PlantStatus,
				}

				siteDocumentBatches = append(siteDocumentBatches, tmp)
			}
		}
	}

	if err := s.elastic.BulkIndex(documentBatches); err != nil {
		s.logger.Error(err)
		return err
	}

	if err := s.elastic.UpsertSiteStation(siteDocumentBatches); err != nil {
		s.logger.Error(err)
		return err
	}

	close(documentCh)
	close(errorCh)
	close(doneCh)

	return nil
}
