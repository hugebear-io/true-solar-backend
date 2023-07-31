package collector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/solarman"
)

type solarmanCollector struct {
	dataCollectorConfig domain.DataCollectorConfigService
	siteRegionMapping   domain.SiteRegionMappingService
	logger              logger.Logger
	now                 time.Time
	beginningOfDay      time.Time
	documentCh          chan interface{}
	errorCh             chan error
	doneCh              chan bool
	brand               string
	usernames           []string
	password            string
	appID               string
	appSecret           string
	siteRegions         []port.SiteRegionMapping
}

func NewSolarmanCollector(
	dataCollectorConfig domain.DataCollectorConfigService,
	siteRegionMapping domain.SiteRegionMappingService,
	logger logger.Logger,
	documentCh chan interface{},
	errorCh chan error,
	doneCh chan bool,
	usernames []string,
	password string,
	siteRegions []port.SiteRegionMapping,
) *solarmanCollector {
	now := time.Now()
	return &solarmanCollector{
		brand:               strings.ToUpper(solarman.BRAND),
		dataCollectorConfig: dataCollectorConfig,
		siteRegionMapping:   siteRegionMapping,
		logger:              logger,
		now:                 now,
		beginningOfDay:      time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location()),
		documentCh:          documentCh,
		errorCh:             errorCh,
		doneCh:              doneCh,
		usernames:           usernames,
		password:            password,
		siteRegions:         siteRegions,
	}
}

func (s solarmanCollector) Run() {
	for _, username := range s.usernames {
		inverter := solarman.NewSolarmanInverter(username, s.password, s.appID, s.appSecret, nil)

		basicTokenResp, err := inverter.GetBasicToken()
		if err != nil {
			s.errorCh <- err
			return
		}

		if helper.EmptyString(basicTokenResp.AccessToken) {
			s.errorCh <- errors.New("access token is empty")
			return
		}

		userInfoResp, err := inverter.GetUserInfo()
		if err != nil {
			s.errorCh <- err
			return
		}

		for _, company := range userInfoResp.OrgInfoList {
			businessTokenResp, err := inverter.GetBusinessToken(company.CompanyID)
			if err != nil {
				s.errorCh <- err
				return
			}

			if helper.EmptyString(businessTokenResp.AccessToken) {
				s.errorCh <- errors.New("access token is empty")
				return
			}

			token := businessTokenResp.AccessToken
			plantList, err := inverter.GetPlantList(token)
			if err != nil {
				s.errorCh <- err
				return
			}

			for _, plant := range plantList {
				stationID := plant.ID
				plantID, _ := helper.ParsePlantID(plant.Name)
				cityName, cityCode, cityArea := helper.ParseSiteID(s.siteRegions, plantID.SiteID)

				plantItem := port.PlantItem{
					Timestamp:         s.now,
					Month:             s.now.Format(constant.MONTH_FORMAT),
					Year:              s.now.Format(constant.YEAR_FORMAT),
					MonthYear:         s.now.Format(constant.MONTH_YEAR_FORMAT),
					VendorType:        s.brand,
					DataType:          constant.DATA_TYPE_PLANT,
					Area:              cityArea,
					SiteID:            plantID.SiteID,
					SiteCityName:      cityName,
					SiteCityCode:      cityCode,
					NodeType:          plantID.NodeType,
					ACPhase:           plantID.ACPhase,
					ID:                strconv.Itoa(stationID),
					Name:              plant.Name,
					Latitude:          plant.LocationLat,
					Longitude:         plant.LocationLng,
					LocationAddress:   plant.LocationAddress,
					InstalledCapacity: plant.InstalledCapacity,
				}

				var mergedElectricPrice float64
				var totalPowerGenerationKWh float64
				var sumYearlyPowerGenerationKWh float64

				if plantItem.Latitude != 0 && plantItem.Longitude != 0 {
					plantItem.Location = fmt.Sprintf("%f,%f", plantItem.Latitude, plantItem.Longitude)
				}

				if plant.CreatedDate != 0 {
					parsed := time.Unix(int64(plant.CreatedDate), 0).UTC()
					plantItem.CreatedDate = &parsed
				}

				if basicInfoResp, err := inverter.GetPlantBaseInfo(token, stationID); err == nil {
					plantItem.Currency = basicInfoResp.Currency
					mergedElectricPrice = basicInfoResp.MergeElectricPrice
				}

				if realtimeDataResp, err := inverter.GetPlantRealtimeData(token, stationID); err == nil {
					plantItem.CurrentPower = realtimeDataResp.GenerationPower / 1000.0
				}

				if resp, err := inverter.GetHistoricalPlantData(
					token,
					stationID,
					solarman.TIME_TYPE_DAY,
					s.now.Unix(),
					s.now.Unix(),
				); err == nil && len(resp.StationDataItems) > 0 {
					plantItem.DailyProduction = resp.StationDataItems[0].GenerationValue
				}

				if resp, err := inverter.GetHistoricalPlantData(
					token,
					stationID,
					solarman.TIME_TYPE_MONTH,
					s.now.Unix(),
					s.now.Unix(),
				); err == nil && len(resp.StationDataItems) > 0 {
					plantItem.MonthlyProduction = resp.StationDataItems[0].GenerationValue
				}

				startTime := time.Date(2015, s.now.Month(), s.now.Day(), 0, 0, 0, 0, s.now.Location())
				if resp, err := inverter.GetHistoricalPlantData(
					token,
					stationID,
					solarman.TIME_TYPE_YEAR,
					startTime.Unix(),
					s.now.Unix(),
				); err == nil && len(resp.StationDataItems) > 0 {
					for _, dataItem := range resp.StationDataItems {
						if dataItem.Year == s.now.Local().Year() {
							plantItem.YearlyProduction = dataItem.GenerationValue
						}

						sumYearlyPowerGenerationKWh += dataItem.GenerationValue
					}
				}

				deviceListResp, err := inverter.GetPlantDeviceList(token, stationID)
				if err != nil {
					s.errorCh <- err
					return
				}

				deviceStatusArray := make([]string, 0)
				for _, device := range deviceListResp {
					deviceSN := device.DeviceSN
					deviceID := device.DeviceID

					deviceItem := port.DeviceItem{
						Timestamp:    s.now,
						Month:        s.now.Format(constant.MONTH_FORMAT),
						Year:         s.now.Format(constant.YEAR_FORMAT),
						MonthYear:    s.now.Format(constant.MONTH_YEAR_FORMAT),
						VendorType:   s.brand,
						DataType:     constant.DATA_TYPE_DEVICE,
						Area:         cityArea,
						SiteID:       plantID.SiteID,
						SiteCityName: cityName,
						SiteCityCode: cityCode,
						NodeType:     plantID.NodeType,
						ACPhase:      plantID.ACPhase,
						PlantID:      plantID.NodeType,
						PlantName:    plant.Name,
						Latitude:     plantItem.Latitude,
						Longitude:    plantItem.Longitude,
						ID:           strconv.Itoa(deviceID),
						SN:           deviceSN,
						Name:         deviceSN,
						DeviceType:   device.DeviceType,
					}

					if resp, err := inverter.GetDeviceRealtimeData(token, deviceSN); err == nil {
						if len(resp.DataList) > 0 {
							for _, data := range resp.DataList {
								if data.Key == solarman.DATA_LIST_KEY_CUMULATIVE_PRODUCTION {
									if generation, err := strconv.ParseFloat(data.Value, 64); err == nil {
										totalPowerGenerationKWh += generation
									}

									break
								}
							}
						}
					}

					if resp, err := inverter.GetHistoricalDeviceData(
						token,
						deviceSN,
						solarman.TIME_TYPE_DAY,
						s.now.Unix(),
						s.now.Unix(),
					); err == nil && len(resp.ParamDataList) > 0 {
						for _, param := range resp.ParamDataList {
							for _, data := range param.DataList {
								if data.Key == solarman.DATA_LIST_KEY_GENERATION {
									if generation, err := strconv.ParseFloat(data.Value, 64); err == nil {
										deviceItem.DailyPowerGeneration += generation
									}

									break
								}
							}
						}
					}

					if resp, err := inverter.GetHistoricalDeviceData(
						token,
						deviceSN,
						solarman.TIME_TYPE_MONTH,
						s.now.Unix(),
						s.now.Unix(),
					); err == nil && len(resp.ParamDataList) > 0 {
						for _, param := range resp.ParamDataList {
							for _, data := range param.DataList {
								if data.Key == solarman.DATA_LIST_KEY_GENERATION {
									if generation, err := strconv.ParseFloat(data.Value, 64); err == nil {
										deviceItem.MonthlyPowerGeneration += generation
									}

									break
								}
							}
						}
					}

					if resp, err := inverter.GetHistoricalDeviceData(
						token,
						deviceSN,
						solarman.TIME_TYPE_YEAR,
						s.now.Unix(),
						s.now.Unix(),
					); err == nil && len(resp.ParamDataList) > 0 {
						for _, param := range resp.ParamDataList {
							for _, data := range param.DataList {
								if data.Key == solarman.DATA_LIST_KEY_GENERATION {
									if generation, err := strconv.ParseFloat(data.Value, 64); err == nil {
										deviceItem.YearlyPowerGeneration += generation
									}

									break
								}
							}
						}

					}

					if device.CollectionTime != 0 {
						parsed := time.Unix(device.CollectionTime, 0).UTC()
						deviceItem.LastUpdateTime = &parsed
					}

					switch device.ConnectStatus {
					case 0:
						deviceItem.Status = solarman.DEVICE_STATUS_OFF
					case 1:
						deviceItem.Status = solarman.DEVICE_STATUS_ON
					case 2:
						deviceItem.Status = solarman.DEVICE_STATUS_FAILURE
						if resp, err := inverter.GetDeviceAlertList(
							token,
							deviceSN,
							s.beginningOfDay.Unix(),
							s.now.Unix(),
						); err == nil {
							for _, alert := range resp {
								alarmItem := port.AlarmItem{
									Timestamp:    s.now,
									Month:        s.now.Format(constant.MONTH_FORMAT),
									Year:         s.now.Format(constant.YEAR_FORMAT),
									MonthYear:    s.now.Format(constant.MONTH_YEAR_FORMAT),
									VendorType:   s.brand,
									DataType:     constant.DATA_TYPE_ALARM,
									Area:         plantID.SiteID,
									SiteID:       plantID.SiteID,
									SiteCityName: cityName,
									SiteCityCode: cityCode,
									NodeType:     plantID.NodeType,
									ACPhase:      plantID.ACPhase,
									PlantID:      strconv.Itoa(stationID),
									PlantName:    plant.Name,
									Latitude:     plantItem.Latitude,
									Longitude:    plantItem.Longitude,
									Location:     plantItem.Location,
									DeviceID:     strconv.Itoa(deviceID),
									DeviceSN:     deviceSN,
									DeviceName:   deviceSN,
									DeviceType:   device.DeviceType,
									DeviceStatus: solarman.DEVICE_STATUS_FAILURE,
									ID:           strconv.Itoa(alert.AlertID),
									Message:      alert.AlertNameInPAAS,
								}

								if alert.AlertTime != 0 {
									alertTime := time.Unix(alert.AlertTime, 0).UTC()
									alarmItem.AlarmTime = &alertTime
								}

								s.documentCh <- alarmItem
							}
						}
					default:
					}

					if helper.EmptyString(deviceItem.Status) {
						deviceStatusArray = append(deviceStatusArray, deviceItem.Status)
					}

					s.documentCh <- deviceItem
				}

				plantStatus := solarman.SOLARMAN_PLANT_STATUS_ON
				if len(deviceStatusArray) > 0 {
					var offlineCount int
					var alertCount int
					for _, status := range deviceStatusArray {
						switch status {
						case solarman.DEVICE_STATUS_OFF:
							offlineCount++
						case solarman.DEVICE_STATUS_ON:
						default:
							alertCount++
						}
					}

					if alertCount > 0 {
						plantStatus = solarman.SOLARMAN_PLANT_STATUS_ALARM
					} else if offlineCount > 0 {
						plantStatus = solarman.SOLARMAN_PLANT_STATUS_OFF
					}
				} else {
					plantStatus = solarman.SOLARMAN_PLANT_STATUS_ON
				}

				plantItem.TotalProduction = totalPowerGenerationKWh
				if plantItem.TotalProduction < plantItem.YearlyProduction {
					plantItem.TotalProduction = sumYearlyPowerGenerationKWh
				}

				plantItem.PlantStatus = plantStatus
				plantItem.TotalSavingPrice = mergedElectricPrice * totalPowerGenerationKWh

				s.documentCh <- plantItem
			}
		}
	}
}
