package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/adapter/collector"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type solarmanCollectorService struct {
	dataCollectorConfig domain.DataCollectorConfigService
	siteRegionConfig    domain.SiteRegionMappingService
	elastic             port.ElasticSearchRepoPort
	logger              logger.Logger
}

func NewSolarmanCollectorService(
	dataCollectorConfig domain.DataCollectorConfigService,
	siteRegionConfig domain.SiteRegionMappingService,
	elastic port.ElasticSearchRepoPort,
	logger logger.Logger,
) domain.SolarmanCollectorService {
	return &solarmanCollectorService{
		dataCollectorConfig: dataCollectorConfig,
		siteRegionConfig:    siteRegionConfig,
		elastic:             elastic,
	}
}

func (s solarmanCollectorService) Run() {
	configs, err := s.dataCollectorConfig.GetDataCollectorConfigByVendorType(constant.VENDOR_TYPE_INVT)
	if err != nil {
		s.logger.Error(err)
		return
	}

	usernames := []string{}
	password := ""
	appID := ""
	appSecret := ""
	for _, config := range configs {
		usernames = append(usernames, config.Username)
		password = config.Password
		appID = *config.AppID
		appSecret = *config.AppSecret
	}

	siteRegions, err := s.siteRegionConfig.GetAllSiteRegionMapping()
	if err != nil {
		s.logger.Error(err)
		return
	}

	usernames = []string{"bignode.invt.th@gmail.com"}
	password = "123456*"
	appID = "202010143565002"
	appSecret = "222c202135013aee622c71cdf8c47757"

	documents := make([]interface{}, 0)
	siteDocuments := make([]port.SiteItem, 0)
	documentCh := make(chan interface{})
	errorCh := make(chan error)
	doneCh := make(chan bool)
	solarmanCollector := collector.NewSolarmanCollector(usernames, password, appID, appSecret, siteRegions, documentCh, errorCh, doneCh)
	go solarmanCollector.Run()

DONE:
	for {
		select {
		case <-doneCh:
			break DONE
		case err := <-errorCh:
			s.logger.Error(err)
		case doc := <-documentCh:
			documents = append(documents, doc)
			if item, ok := doc.(port.PlantItem); ok {
				tmp := port.SiteItem{
					Timestamp:   &item.Timestamp,
					VendorType:  item.VendorType,
					Area:        item.Area,
					SiteID:      item.SiteID,
					NodeType:    item.NodeType,
					Name:        item.Name,
					Location:    item.Location,
					PlantStatus: item.PlantStatus,
				}

				siteDocuments = append(siteDocuments, tmp)
			}
		}
	}

	if err := s.elastic.BulkIndex(documents); err != nil {
		s.logger.Error(err)
		return
	}

	if err := s.elastic.UpsertSiteStation(siteDocuments); err != nil {
		s.logger.Error(err)
		return
	}
}
