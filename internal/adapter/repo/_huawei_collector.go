package repo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/huawei"
)

type huaweiCollectorRepo struct {
	now         time.Time
	siteRegions []port.SiteRegionMapping
	logger      logger.Logger
}

func NewHuaweiCollectorRepo(siteRegions []port.SiteRegionMapping, logger logger.Logger) *huaweiCollectorRepo {
	obj := huaweiCollectorRepo{}
	obj.siteRegions = siteRegions
	obj.logger = logger
	obj.now = time.Now()
	return &obj
}

func (r huaweiCollectorRepo) CollectData(username, password string, documentCh chan interface{}, doneCh chan bool, errorCh chan error) {
	beginTime := time.Date(r.now.Year(), r.now.Month(), r.now.Day(), 6, 0, 0, 0, r.now.Location()).Unix()
	collectTime := r.now.Unix()

	inverter := huawei.NewHuaweiInverter(username, password, nil)
	plantList, err := inverter.GetPlantList()
	if err != nil {
		errorCh <- err
		return
	}

	var arrPlantCode []string
	var arrPlantCodeStr []string
	for _, plant := range plantList {
		if len(arrPlantCode) == 100 {
			arrPlantCodeStr = append(arrPlantCodeStr, strings.Join(arrPlantCode, ","))
		}

		if !helper.EmptyString(plant.Code) {
			arrPlantCode = append(arrPlantCode, plant.Code)
		}
	}
	arrPlantCodeStr = append(arrPlantCodeStr, strings.Join(arrPlantCode, ","))

	var arrDeviceItem []huawei.DeviceItem
	mapPlantCodeToRealtimeData := make(map[string]huawei.PlantRealtimeData)
	mapPlantCodeToDailyData := make(map[string]huawei.HistoricalPlantData)
	mapPlantCodeToMonthlyData := make(map[string]huawei.HistoricalPlantData)
	mapPlantCodeToYearlyPower := make(map[string]float64)
	mapPlantCodeToTotalPower := make(map[string]float64)
	mapPlantCodeToTotalCO2 := make(map[string]float64)
	mapPlantCodeToDevice := make(map[string][]huawei.DeviceItem)
	mapDeviceSNToAlarm := make(map[string][]huawei.DeviceAlarmItem)

	for _, codeStr := range arrPlantCodeStr {
		plantRealtimeData, err := inverter.GetPlantRealtimeData(codeStr)
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range plantRealtimeData {
			if !helper.EmptyString(item.Code) {
				mapPlantCodeToRealtimeData[item.Code] = item
			}
		}

		plantDailyData, err := inverter.GetDailyPlantData(codeStr, r.now.Unix())
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range plantDailyData {
			if !helper.EmptyString(item.Code) {
				if r.now.Format("2006-01-02") == time.Unix(item.CollectTime/1e3, 0).Format("2006-01-02") {
					mapPlantCodeToDailyData[item.Code] = item
				}
			}
		}

		plantMonthlyData, err := inverter.GetMonthlyPlantData(codeStr, r.now.Unix())
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range plantMonthlyData {
			if !helper.EmptyString(item.Code) {
				if r.now.Format("2006-01") == time.Unix(item.CollectTime/1e3, 0).Format("2006-01") {
					mapPlantCodeToMonthlyData[item.Code] = item
				}

				mapPlantCodeToYearlyPower[item.Code] = mapPlantCodeToYearlyPower[item.Code] + item.DataItemMap.InverterPower
			}
		}

		plantYearlyData, err := inverter.GetYearlyPlantData(codeStr, r.now.Unix())
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range plantYearlyData {
			if !helper.EmptyString(item.Code) {
				mapPlantCodeToTotalPower[item.Code] = mapPlantCodeToTotalPower[item.Code] + item.DataItemMap.InverterPower
				mapPlantCodeToTotalCO2[item.Code] = mapPlantCodeToTotalCO2[item.Code] + item.DataItemMap.ReductionTotalCO2
			}
		}

		deviceList, err := inverter.GetDeviceList(codeStr)
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range deviceList {
			if !helper.EmptyString(item.PlantCode) {
				mapPlantCodeToDevice[item.PlantCode] = append(mapPlantCodeToDevice[item.PlantCode], item)
			}

			if item.TypeID == 1 {
				arrDeviceItem = append(arrDeviceItem, item)
			}
		}

		deviceAlarmList, err := inverter.GetDeviceAlarmList(codeStr, beginTime, collectTime)
		if err != nil {
			errorCh <- err
			return
		}

		for _, alarm := range deviceAlarmList {
			doubleAlarm := false

			if !helper.EmptyString(alarm.SerialNumber) {
				for i, item := range mapDeviceSNToAlarm[alarm.SerialNumber] {
					if item.AlarmName == alarm.AlarmName {
						doubleAlarm = true
						if item.RaiseTime < alarm.RaiseTime {
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

	var arrDeviceID []string
	var arrDeviceIDString []string
	for _, item := range arrDeviceItem {
		if len(arrDeviceID) == 100 {
			tmp := strings.Join(arrDeviceID, ",")
			arrDeviceIDString = append(arrDeviceIDString, tmp)
			arrDeviceID = []string{}
		}

		if item.ID != 0 {
			idString := strconv.Itoa(item.ID)
			arrDeviceID = append(arrDeviceID, idString)
		}
	}
	arrDeviceIDString = append(arrDeviceIDString, strings.Join(arrDeviceID, ","))

	mapDeviceToRealtimeData := make(map[int]huawei.RealtimeDeviceData)
	mapDeviceToDailyData := make(map[int]huawei.HistoricalDeviceData)
	mapDeviceToMonthlyData := make(map[int]huawei.HistoricalDeviceData)
	mapDeviceToYearlyPower := make(map[int]float64)
	for _, deviceID := range arrDeviceIDString {
		realtimeDeviceData, err := inverter.GetDeviceRealtimeData(deviceID, "1")
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range realtimeDeviceData {
			if item.ID != 0 {
				mapDeviceToRealtimeData[item.ID] = item
			}
		}

		dailyDeviceData, err := inverter.GetDailyDeviceData(deviceID, "1", collectTime)
		if err != nil {
			errorCh <- err
			return
		}

		for _, item := range dailyDeviceData {
			if item.ID != 0 {
				if r.now.Format("2006-01-02") == time.Unix(item.CollectTime/1e3, 0).Format("2006-01-02") {
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
			errorCh <- err
			return
		}

		for _, item := range monthlyDeviceData {
			if item.ID != nil {
				deviceID := item.ID
				switch deviceID := deviceID.(type) {
				case float64:
					parsedDeviceID := int(deviceID)
					mapDeviceToYearlyPower[parsedDeviceID] = mapDeviceToYearlyPower[parsedDeviceID] + item.DataItemMap.ProductPower
					if r.now.Format("2006-01") == time.Unix(item.CollectTime/1e3, 0).Format("2006-01") {
						mapDeviceToMonthlyData[parsedDeviceID] = item
					}
				default:
				}
			}
		}
	}

	for _, plant := range plantList {
		plantNameInfo, _ := helper.ParsePlantID(plant.Name)
		cityName, cityCode, cityArea := helper.ParseSiteID(r.siteRegions, plantNameInfo.SiteID)

		var latitude float64
		var longitude float64
		var location string
		currentPower := 0.0
		plantStatus := huawei.HuaweiMapDeviceStatus[mapPlantCodeToRealtimeData[plant.Code].DataItemMap.RealHealthState]

		for _, device := range mapPlantCodeToDevice[plant.Code] {
			if device.Latitude != 0 && device.Longitude != 0 {
				location = fmt.Sprintf("%f,%f", device.Latitude, device.Longitude)
			}

			var deviceStatus int
			if mapDeviceToRealtimeData[device.ID].DataItemMap.Status != 0 {
				deviceStatus = mapDeviceToRealtimeData[device.ID].DataItemMap.Status
			}

			if len(mapDeviceSNToAlarm[device.SerialNumber]) > 0 {
				deviceStatus = 2
				for _, deviceAlarm := range mapDeviceSNToAlarm[device.SerialNumber] {
					alarmItem := port.AlarmItem{
						Timestamp:    r.now,
						Month:        r.now.Format("01"),
						Year:         r.now.Format("2006"),
						MonthYear:    r.now.Format("01-2006"),
						VendorType:   huawei.BRAND,
						DataType:     constant.DATA_TYPE_ALARM,
						Area:         cityArea,
						SiteID:       plantNameInfo.SiteID,
						SiteCityCode: cityCode,
						SiteCityName: cityName,
						NodeType:     plantNameInfo.SiteID,
						ACPhase:      plantNameInfo.ACPhase,
						PlantID:      plant.Code,
						PlantName:    plant.Name,
						Latitude:     device.Latitude,
						Longitude:    device.Longitude,
						Location:     location,
						DeviceID:     strconv.Itoa(device.ID),
						DeviceSN:     device.SerialNumber,
						DeviceName:   device.Name,
						DeviceStatus: huawei.HUAWEI_STATUS_ALARM,
						ID:           strconv.Itoa(deviceAlarm.AlarmID),
						Message:      deviceAlarm.AlarmName,
					}

					if deviceAlarm.RaiseTime != 0 {
						timeUTC := time.Unix(deviceAlarm.RaiseTime/1e3, 0).UTC()
						alarmItem.AlarmTime = &timeUTC
					}

					documentCh <- alarmItem
				}
			}

			deviceItem := port.DeviceItem{
				Timestamp:      r.now,
				Month:          r.now.Format("01"),
				Year:           r.now.Format("2006"),
				MonthYear:      r.now.Format("01-2006"),
				VendorType:     huawei.BRAND,
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

			documentCh <- deviceItem
		}

		plantItem := port.PlantItem{
			Timestamp:         r.now,
			Month:             r.now.Format("01"),
			Year:              r.now.Format("2006"),
			MonthYear:         r.now.Format("01-2006"),
			VendorType:        huawei.BRAND,
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
			DailyProduction:   mapPlantCodeToDailyData[plant.Code].DataItemMap.InverterPower,
			MonthlyProduction: mapPlantCodeToMonthlyData[plant.Code].DataItemMap.InverterPower,
			YearlyProduction:  mapPlantCodeToYearlyPower[plant.Code],
			PlantStatus:       plantStatus,
		}

		plantItem.TotalProduction = mapPlantCodeToRealtimeData[plant.Code].DataItemMap.TotalPower
		if plantItem.TotalProduction < plantItem.YearlyProduction {
			plantItem.TotalProduction = mapPlantCodeToTotalPower[plant.Code]
		}

		documentCh <- plantItem
	}

	doneCh <- true
}
