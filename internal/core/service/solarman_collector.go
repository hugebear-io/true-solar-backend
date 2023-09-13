package service

import (
	"fmt"

	"github.com/hugebear-io/true-solar-backend/internal/adapter/collector"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type solarmanCollectorService struct {
	dataCollectorConfig domain.DataCollectorConfigService
	siteRegionConfig    domain.SiteRegionMappingService
	logger              logger.Logger
}

// TODO: REFACTOR
func NewSolarmanCollectorService(
	dataCollectorConfig domain.DataCollectorConfigService,
	siteRegionConfig domain.SiteRegionMappingService,
) domain.SolarmanCollectorService {
	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/solarman-collector-service.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(logger.LOG_LEVEL_DEBUG),
		SkipCaller:  1,
	})

	return &solarmanCollectorService{
		dataCollectorConfig: dataCollectorConfig,
		siteRegionConfig:    siteRegionConfig,
		logger:              l,
	}
}

func (s solarmanCollectorService) Run() {
	defer s.Close()

	elasticConfig := config.Config.ElasticSearch
	elasticClient := infra.NewElasticSearch(s.logger)
	elastic := repo.NewElasticSearchRepo(elasticClient, elasticConfig.Index)

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

	// usernames = []string{"bignode.invt.th@gmail.com"}
	// password = "123456*"
	// appID = "202010143565002"
	// appSecret = "222c202135013aee622c71cdf8c47757"

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
			fmt.Printf("%#v\n", doc)

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

	if err := elastic.BulkIndex(documents); err != nil {
		s.logger.Error(err)
		return
	}

	if err := elastic.UpsertSiteStation(siteDocuments); err != nil {
		s.logger.Error(err)
		return
	}
}

func (s *solarmanCollectorService) Close() {
	s.logger.Close()
}
