package collector

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/huawei"
)

const (
	dateFormat      = "2006-01-02"
	monthYearFormat = "2006-01"
	yearFormat      = "2006"
	monthFormat     = "01"
)

type huaweiCollector struct {
	dataCollectorConfig domain.DataCollectorConfigService
	siteRegionMapping   domain.SiteRegionMappingService
	logger              logger.Logger
	now                 time.Time
	documentCh          chan interface{}
	errorCh             chan error
	doneCh              chan bool
	brand               string
	usernames           []string
	password            string
	siteRegions         []port.SiteRegionMapping
}

func NewHuaweiCollector(
	dataCollectorConfig domain.DataCollectorConfigService,
	siteRegionMapping domain.SiteRegionMappingService,
	logger logger.Logger,
	documentCh chan interface{},
	errorCh chan error,
	doneCh chan bool,
	usernames []string,
	password string,
	siteRegions []port.SiteRegionMapping,
) port.HuaweiCollector {
	now := time.Now()
	return &huaweiCollector{
		dataCollectorConfig: dataCollectorConfig,
		siteRegionMapping:   siteRegionMapping,
		logger:              logger,
		brand:               strings.ToUpper(huawei.BRAND),
		now:                 now,
		documentCh:          documentCh,
		errorCh:             errorCh,
		doneCh:              doneCh,
		usernames:           usernames,
		password:            password,
		siteRegions:         siteRegions,
	}
}

func (c huaweiCollector) Run() {
	beginTime := time.Date(c.now.Year(), c.now.Month(), c.now.Day(), 6, 0, 0, 0, c.now.Location()).UnixNano() / 1e6
	collectTime := c.now.UnixNano() / 1e6

	for _, username := range c.usernames {
		inverter := huawei.NewHuaweiInverter(username, c.password, nil)

		plantList, err := inverter.GetPlantList()
		if err != nil {
			c.errorCh <- err
			return
		}

		var plantCodeList []string
		var plantCodeListString []string
		for _, plant := range plantList {
			if len(plantCodeList) == 100 {
				tmp := strings.Join(plantCodeList, ",")
				plantCodeListString = append(plantCodeListString, tmp)
				plantCodeList = []string{}
			}

			if !helper.EmptyString(plant.Code) {
				plantCodeList = append(plantCodeList, plant.Code)
			}
		}
		plantCodeListString = append(plantCodeListString, strings.Join(plantCodeList, ","))

		var deviceList []huawei.DeviceItem
		mapPlantCodeToRealtimeData := make(map[string]huawei.PlantRealtimeData)
		mapPlantCodeToDailyData := make(map[string]huawei.HistoricalPlantData)
		mapPlantCodeToMonthlyData := make(map[string]huawei.HistoricalPlantData)
		mapPlantCodeToYearlyPower := make(map[string]float64)
		mapPlantCodeToTotalPower := make(map[string]float64)
		mapPlantCodeToTotalCO2 := make(map[string]float64)
		mapPlantCodeToDevice := make(map[string][]huawei.DeviceItem)
		mapDeviceSNToAlarm := make(map[string][]huawei.DeviceAlarmItem)

		for _, plantCode := range plantCodeListString {
			realtimePlantData, err := inverter.GetPlantRealtimeData(plantCode)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range realtimePlantData {
				if !helper.EmptyString(item.Code) {
					mapPlantCodeToRealtimeData[item.Code] = item
				}
			}

			dailyPlantData, err := inverter.GetDailyPlantData(plantCode, collectTime)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range dailyPlantData {
				if !helper.EmptyString(item.Code) {
					if c.now.Format(dateFormat) == time.Unix(item.CollectTime/1e3, 0).Format(dateFormat) {
						mapPlantCodeToDailyData[item.Code] = item
					}
				}
			}

			monthlyPlantData, err := inverter.GetMonthlyPlantData(plantCode, collectTime)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range monthlyPlantData {
				if !helper.EmptyString(item.Code) {
					if c.now.Format(monthYearFormat) == time.Unix(item.CollectTime/1e3, 0).Format(monthYearFormat) {
						mapPlantCodeToMonthlyData[item.Code] = item
					}

					mapPlantCodeToYearlyPower[item.Code] = mapPlantCodeToYearlyPower[item.Code] + item.DataItemMap.InverterPower
				}
			}

			yearlyPlantData, err := inverter.GetYearlyPlantData(plantCode, collectTime)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range yearlyPlantData {
				if !helper.EmptyString(item.Code) {
					mapPlantCodeToTotalPower[item.Code] = mapPlantCodeToTotalPower[item.Code] + item.DataItemMap.InverterPower
					mapPlantCodeToTotalCO2[item.Code] = mapPlantCodeToTotalCO2[item.Code] + item.DataItemMap.ReductionTotalCO2
				}
			}

			deviceListResp, err := inverter.GetDeviceList(plantCode)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range deviceListResp {
				if !helper.EmptyString(item.PlantCode) {
					mapPlantCodeToDevice[item.PlantCode] = append(mapPlantCodeToDevice[item.PlantCode], item)
				}

				if item.TypeID == 1 {
					deviceList = append(deviceList, item)
				}
			}

			deviceAlarmList, err := inverter.GetDeviceAlarmList(plantCode, beginTime, collectTime)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, alarm := range deviceAlarmList {
				doubleAlarm := false

				if !helper.EmptyString(alarm.SerialNumber) {
					for i, deviceAlarm := range mapDeviceSNToAlarm[alarm.SerialNumber] {
						if deviceAlarm.AlarmName == alarm.AlarmName {
							doubleAlarm = true

							if deviceAlarm.RaiseTime < alarm.RaiseTime {
								mapDeviceSNToAlarm[alarm.SerialNumber][i] = alarm
								break
							}
						}
					}

					if !doubleAlarm {
						mapDeviceSNToAlarm[alarm.SerialNumber] = append(mapDeviceSNToAlarm[alarm.SerialNumber], alarm)
					}
				}
			}
		}

		var deviceIDList []string
		var deviceIDListString []string
		for _, item := range deviceList {
			if len(deviceIDList) == 100 {
				tmp := strings.Join(deviceIDList, ",")
				deviceIDListString = append(deviceIDListString, tmp)
				deviceIDList = []string{}
			}

			if item.ID != 0 {
				deviceIDList = append(deviceIDList, strconv.Itoa(item.ID))
			}
		}
		deviceIDListString = append(deviceIDListString, strings.Join(deviceIDList, ","))

		mapDeviceToRealtimeData := make(map[int]huawei.RealtimeDeviceData)
		mapDeviceToDailyData := make(map[int]huawei.HistoricalDeviceData)
		mapDeviceToMonthlyData := make(map[int]huawei.HistoricalDeviceData)
		mapDeviceToYearlyPower := make(map[int]float64)

		for _, deviceID := range deviceIDListString {
			realtimeDeviceData, err := inverter.GetDeviceRealtimeData(deviceID, "1")
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range realtimeDeviceData {
				if item.ID != 0 {
					mapDeviceToRealtimeData[item.ID] = item
				}
			}

			dailyDeviceData, err := inverter.GetDailyDeviceData(deviceID, "1", collectTime)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range dailyDeviceData {
				if item.ID != 0 {
					if c.now.Format(dateFormat) == time.Unix(item.CollectTime/1e3, 0).Format(dateFormat) {
						deviceID := item.ID
						switch deviceID := deviceID.(type) {
						case float64:
							parsedDeviceID := int(deviceID)
							mapDeviceToDailyData[parsedDeviceID] = item
						default:
						}
					}
				}
			}

			monthlyDeviceData, err := inverter.GetMonthlyDeviceData(deviceID, "1", collectTime)
			if err != nil {
				c.errorCh <- err
				return
			}

			for _, item := range monthlyDeviceData {
				if item.ID != nil {
					switch id := item.ID.(type) {
					case float64:
						parsedID := int(id)
						mapDeviceToYearlyPower[parsedID] = mapDeviceToYearlyPower[parsedID] + item.DataItemMap.ProductPower

						if c.now.Format(monthYearFormat) == time.Unix(item.CollectTime/1e3, 0).Format(monthYearFormat) {
							mapDeviceToMonthlyData[parsedID] = item
						}
					default:
					}
				}
			}
		}

		for _, plant := range plantList {
			plantNameInfo, _ := helper.ParsePlantID(plant.Name)
			cityName, cityCode, cityArea := helper.ParseSiteID(c.siteRegions, plantNameInfo.SiteID)

			var latitude, longitude float64
			var location string
			currentPower := 0.0
			plantStatus := huawei.HuaweiMapPlantStatus[mapPlantCodeToRealtimeData[plant.Code].DataItemMap.RealHealthState]

			for _, device := range mapPlantCodeToDevice[plant.Code] {
				if device.Latitude != 0 && device.Longitude != 0 {
					location = fmt.Sprintf("%v,%v", device.Longitude, device.Latitude)
				}

				var deviceStatus int
				if mapDeviceToRealtimeData[device.ID].DataItemMap.Status != 0 {
					deviceStatus = mapDeviceToRealtimeData[device.ID].DataItemMap.Status
				}

				if len(mapDeviceSNToAlarm[device.SerialNumber]) > 0 {
					deviceStatus = 2

					for _, deviceAlarm := range mapDeviceSNToAlarm[device.SerialNumber] {
						alarmItem := port.AlarmItem{
							Timestamp:    c.now,
							Month:        c.now.Format(monthFormat),
							Year:         c.now.Format(yearFormat),
							MonthYear:    c.now.Format(monthYearFormat),
							VendorType:   c.brand,
							DataType:     constant.DATA_TYPE_ALARM,
							Area:         cityArea,
							SiteID:       plantNameInfo.SiteID,
							SiteCityCode: cityCode,
							SiteCityName: cityName,
							NodeType:     plantNameInfo.NodeType,
							ACPhase:      plantNameInfo.ACPhase,
							PlantID:      plant.Code,
							PlantName:    plant.Name,
							Latitude:     device.Latitude,
							Longitude:    device.Longitude,
							Location:     location,
							DeviceID:     strconv.Itoa(device.ID),
							DeviceSN:     deviceAlarm.SerialNumber,
							DeviceName:   deviceAlarm.DeviceName,
							DeviceStatus: huawei.HUAWEI_STATUS_ALARM,
							ID:           strconv.Itoa(deviceAlarm.AlarmID),
							Message:      deviceAlarm.AlarmName,
						}

						if deviceAlarm.RaiseTime != 0 {
							timeUTC := time.Unix(deviceAlarm.RaiseTime/1e3, 0).UTC()
							alarmItem.AlarmTime = &timeUTC
						}

						c.documentCh <- alarmItem
					}
				}

				deviceItem := port.DeviceItem{
					Timestamp:      c.now,
					Month:          c.now.Format(monthFormat),
					Year:           c.now.Format(yearFormat),
					MonthYear:      c.now.Format(monthYearFormat),
					VendorType:     c.brand,
					DataType:       constant.DATA_TYPE_DEVICE,
					Area:           cityArea,
					SiteID:         plantNameInfo.SiteID,
					SiteCityCode:   cityCode,
					SiteCityName:   cityName,
					NodeType:       plantNameInfo.NodeType,
					ACPhase:        plantNameInfo.ACPhase,
					PlantID:        plant.Code,
					PlantName:      plant.Name,
					Latitude:       device.Latitude,
					Longitude:      device.Longitude,
					Location:       location,
					ID:             strconv.Itoa(device.ID),
					SN:             device.SerialNumber,
					Name:           device.Name,
					LastUpdateTime: nil,
				}

				switch deviceStatus {
				case 0:
					deviceItem.Status = huawei.HUAWEI_STATUS_OFF
					if plantStatus != huawei.HUAWEI_STATUS_ALARM {
						plantStatus = huawei.HUAWEI_STATUS_OFF
					}
				case 1:
					deviceItem.Status = huawei.HUAWEI_STATUS_ON
					if plantStatus != huawei.HUAWEI_STATUS_OFF && plantStatus != huawei.HUAWEI_STATUS_ALARM {
						plantStatus = huawei.HUAWEI_STATUS_ON
					}
				case 2:
					deviceItem.Status = huawei.HUAWEI_STATUS_ALARM
					plantStatus = huawei.HUAWEI_STATUS_ALARM
				default:
				}

				if device.TypeID != 0 {
					deviceItem.DeviceType = huawei.HuaweiMapDeviceType[device.TypeID]

					if device.TypeID == 1 {
						deviceItem.TotalPowerGeneration = mapDeviceToRealtimeData[device.ID].DataItemMap.TotalEnergy
						deviceItem.DailyPowerGeneration = mapDeviceToDailyData[device.ID].DataItemMap.ProductPower
						deviceItem.MonthlyPowerGeneration = mapDeviceToMonthlyData[device.ID].DataItemMap.ProductPower
						deviceItem.YearlyPowerGeneration = mapDeviceToYearlyPower[device.ID]

						currentPower += mapDeviceToRealtimeData[device.ID].DataItemMap.ActivePower
						latitude = deviceItem.Latitude
						longitude = deviceItem.Longitude
					}
				}

				c.documentCh <- deviceItem
			}

			plantItem := port.PlantItem{
				Timestamp:         c.now,
				Month:             c.now.Format(monthFormat),
				Year:              c.now.Format(yearFormat),
				MonthYear:         c.now.Format(monthYearFormat),
				VendorType:        c.brand,
				DataType:          constant.DATA_TYPE_PLANT,
				Area:              cityArea,
				SiteID:            plantNameInfo.SiteID,
				SiteCityCode:      cityCode,
				SiteCityName:      cityName,
				NodeType:          plantNameInfo.NodeType,
				ACPhase:           plantNameInfo.ACPhase,
				ID:                plant.Code,
				Name:              plant.Name,
				Latitude:          latitude,
				Longitude:         longitude,
				Location:          location,
				LocationAddress:   plant.Address,
				CreatedDate:       nil,
				InstalledCapacity: plant.Capacity * 1000,
				TotalCO2:          mapPlantCodeToTotalCO2[plant.Code] * 1000,
				MonthlyCO2:        mapPlantCodeToMonthlyData[plant.Code].DataItemMap.ReductionTotalCO2 * 1000,
				TotalSavingPrice:  mapPlantCodeToRealtimeData[plant.Code].DataItemMap.TotalIncome,
				Currency:          huawei.HUAWEI_CURRENCY_USD,
				CurrentPower:      currentPower,
				DailyProduction:   mapPlantCodeToDailyData[plant.Code].DataItemMap.InverterPower,
				MonthlyProduction: mapPlantCodeToMonthlyData[plant.Code].DataItemMap.InverterPower,
				YearlyProduction:  mapPlantCodeToYearlyPower[plant.Code],
				PlantStatus:       plantStatus,
			}

			plantItem.TotalProduction = mapPlantCodeToRealtimeData[plant.Code].DataItemMap.TotalPower
			if plantItem.TotalProduction < plantItem.YearlyProduction {
				plantItem.TotalProduction = mapPlantCodeToTotalPower[plant.Code]
			}

			c.documentCh <- plantItem
		}
	}

	c.doneCh <- true
}
