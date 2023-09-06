package repo

import (
	"fmt"
	"time"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/kstar"
)

type kstarCollectorRepo struct {
	now         time.Time
	siteRegions []port.SiteRegionMapping
	logger      logger.Logger
}

func NewKstarCollectorRepo(siteRegions []port.SiteRegionMapping, logger logger.Logger) *kstarCollectorRepo {
	return &kstarCollectorRepo{now: time.Now(), siteRegions: siteRegions, logger: logger}
}

func (r kstarCollectorRepo) CollectData(username, password string, documentCh chan interface{}, doneCh chan bool, errorCh chan error) {
	inverter := kstar.NewKStarInverter(username, password, &kstar.KStarInverterOptions{})

	mapPlantIDToDeviceList := make(map[string][]kstar.DeviceItem)
	deviceList, err := inverter.GetDeviceList()
	if err != nil {
		errorCh <- err
		return
	}

	for _, device := range deviceList {
		mapPlantIDToDeviceList[device.PlantID] = append(mapPlantIDToDeviceList[device.PlantID], device)
	}

	plantList, err := inverter.GetPlantList()
	if err != nil {
		errorCh <- err
		return
	}

	for _, plant := range plantList {
		plantNameInfo, _ := helper.ParsePlantID(plant.ID)
		cityName, cityCode, cityArea := helper.ParseSiteID(r.siteRegions, plantNameInfo.SiteID)

		var plantStatus string
		var currentPower float64
		var totalProduction float64
		var dailyProduction float64
		var monthlyProduction float64
		var yearlyProduction float64
		var location string

		if plant.Latitude != 0 && plant.Longitude != 0 {
			location = fmt.Sprintf("%f,%f", plant.Latitude, plant.Longitude)
		}

		for _, device := range mapPlantIDToDeviceList[plant.ID] {
			realtimeAlarmList, err := inverter.GetRealtimeAlarmListOfDevice(device.ID)
			if err != nil {
				errorCh <- err
				return
			}

			deviceStatus := device.Status
			if len(realtimeAlarmList) > 0 {
				deviceStatus = 2

				for _, alarm := range realtimeAlarmList {
					alarmItem := port.AlarmItem{
						Timestamp:    r.now,
						Month:        r.now.Format("01"),
						Year:         r.now.Format("2006"),
						MonthYear:    r.now.Format("01-2006"),
						VendorType:   kstar.BRAND,
						DataType:     constant.DATA_TYPE_ALARM,
						Area:         cityArea,
						SiteID:       plantNameInfo.SiteID,
						SiteCityName: cityName,
						SiteCityCode: cityCode,
						NodeType:     plantNameInfo.NodeType,
						ACPhase:      plantNameInfo.ACPhase,
						PlantID:      alarm.PlantID,
						PlantName:    alarm.PlantName,
						Latitude:     plant.Latitude,
						Longitude:    plant.Longitude,
						Location:     location,
						DeviceID:     alarm.DeviceID,
						DeviceSN:     device.InverterID,
						DeviceName:   device.Name,
						DeviceStatus: kstar.KSTAR_DEVICE_STATUS_ALARM,
						ID:           "",
						Message:      alarm.Message,
					}

					if !helper.EmptyString(alarm.SaveTime) {
						if alarmTime, err := time.Parse("2006-01-02 15:04:05", alarm.SaveTime); err != nil {
							timeUTC := alarmTime.UTC()
							alarmItem.AlarmTime = &timeUTC
						}
					}

					documentCh <- alarmItem
				}
			}

			deviceItem := port.DeviceItem{
				Timestamp:    r.now,
				Month:        r.now.Format("01"),
				Year:         r.now.Format("2006"),
				MonthYear:    r.now.Format("01-2006"),
				VendorType:   kstar.BRAND,
				DataType:     constant.DATA_TYPE_DEVICE,
				Area:         cityArea,
				SiteID:       plantNameInfo.SiteID,
				SiteCityCode: cityCode,
				SiteCityName: cityName,
				NodeType:     plantNameInfo.NodeType,
				ACPhase:      plantNameInfo.ACPhase,
				PlantID:      device.PlantID,
				PlantName:    device.PlantName,
				Latitude:     plant.Latitude,
				Longitude:    plant.Longitude,
				Location:     location,
				ID:           device.ID,
				SN:           device.InverterID,
				Name:         device.Name,
				DeviceType:   kstar.KSTAR_DEVICE_TYPE_INVERTER,
			}

			deviceInfo, err := inverter.GetRealtimeDeviceData(device.ID)
			if err != nil {
				errorCh <- err
				return
			}

			if deviceInfo != nil {
				if !helper.EmptyString(device.SaveTime) {
					if saveTime, err := time.Parse("2006-01-02 15:04:05", deviceInfo.SaveTime); err != nil {
						timeUTC := saveTime.UTC()
						deviceItem.LastUpdateTime = &timeUTC
					}
				}

				deviceItem.TotalPowerGeneration = deviceInfo.TotalGeneration
				deviceItem.DailyPowerGeneration = deviceInfo.DayGeneration
				deviceItem.MonthlyPowerGeneration = deviceInfo.MonthGeneration
				deviceItem.YearlyPowerGeneration = deviceInfo.YearGeneration

				currentPower += deviceInfo.PowerInter
				totalProduction += deviceInfo.TotalGeneration
				dailyProduction += deviceInfo.DayGeneration
				monthlyProduction += deviceInfo.MonthGeneration
				yearlyProduction += deviceInfo.YearGeneration

				switch deviceStatus {
				case 0:
					deviceItem.Status = kstar.KSTAR_DEVICE_STATUS_OFF
					if plantStatus != kstar.KSTAR_DEVICE_STATUS_ALARM {
						plantStatus = kstar.KSTAR_DEVICE_STATUS_OFF
					}
				case 1:
					deviceItem.Status = kstar.KSTAR_DEVICE_STATUS_ON
					if plantStatus != kstar.KSTAR_DEVICE_STATUS_ALARM {
						plantStatus = kstar.KSTAR_DEVICE_STATUS_ON
					}
				case 2:
					deviceItem.Status = kstar.KSTAR_DEVICE_STATUS_ALARM
					plantStatus = kstar.KSTAR_DEVICE_STATUS_ALARM
				default:
				}
			}

			documentCh <- deviceItem
		}

		if !helper.EmptyString(plantStatus) {
			plantStatus = kstar.KSTAR_DEVICE_STATUS_OFF
		}

		plantItem := port.PlantItem{
			Timestamp:         r.now,
			Month:             r.now.Format("01"),
			Year:              r.now.Format("2006"),
			MonthYear:         r.now.Format("01-2006"),
			VendorType:        kstar.BRAND,
			DataType:          constant.DATA_TYPE_PLANT,
			Area:              cityArea,
			SiteID:            plantNameInfo.SiteID,
			SiteCityName:      cityName,
			SiteCityCode:      cityCode,
			NodeType:          plantNameInfo.NodeType,
			ACPhase:           plantNameInfo.ACPhase,
			ID:                plant.ID,
			Name:              plant.Name,
			Latitude:          plant.Latitude,
			Longitude:         plant.Longitude,
			Location:          location,
			LocationAddress:   plant.Address,
			InstalledCapacity: plant.InstalledCapacity,
			TotalCO2:          0,
			MonthlyCO2:        0,
			TotalSavingPrice:  totalProduction * plant.ElectricPrice,
			Currency:          plant.ElectricUnit,
			CurrentPower:      currentPower / 1000,
			TotalProduction:   totalProduction,
			DailyProduction:   dailyProduction,
			MonthlyProduction: monthlyProduction,
			YearlyProduction:  yearlyProduction,
			PlantStatus:       plantStatus,
		}

		if !helper.EmptyString(plant.CreatedTime) {
			if createdTime, err := time.Parse("2006-01-02T15:04:05.000+0000", plant.CreatedTime); err != nil {
				timeUTC := createdTime.UTC()
				plantItem.CreatedDate = &timeUTC
			}
		}

		documentCh <- plantItem
	}

	doneCh <- true
}
