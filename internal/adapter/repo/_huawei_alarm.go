package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/huawei"
)

type huaweiAlarmRepo struct {
	rdb    *redis.Client
	snmp   port.SNMPRepoPort
	logger logger.Logger
}

func NewHuaweiAlarmRepo(rdb *redis.Client, snmp port.SNMPRepoPort, logger logger.Logger) *huaweiAlarmRepo {
	return &huaweiAlarmRepo{
		rdb:    rdb,
		snmp:   snmp,
		logger: logger,
	}
}

func (r huaweiAlarmRepo) Alarm(username, password string, now *time.Time) error {
	ctx := context.Background()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location()).UnixNano() / 1e6
	endTime := now.UnixNano() / 1e6
	inverter := huawei.NewHuaweiInverter(username, password, nil)

	plantList, err := inverter.GetPlantList()
	if err != nil {
		r.logger.Error(err)
		return err
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

	var arrDeviceItem []huawei.DeviceItem
	mapPlantCodeToDevice := make(map[string][]huawei.DeviceItem)
	mapDeviceSNToAlarm := make(map[string][]huawei.DeviceAlarmItem)
	mapInverterIDToRealtimeData := make(map[int]huawei.RealtimeDeviceData)

	for _, plantCode := range plantCodeListString {
		deviceList, err := inverter.GetDeviceList(plantCode)
		if err != nil {
			r.logger.Error(err)
			return err
		}

		for _, device := range deviceList {
			if !helper.EmptyString(device.PlantCode) {
				mapPlantCodeToDevice[device.PlantCode] = append(mapPlantCodeToDevice[device.PlantCode], device)
			}

			if device.TypeID == 1 {
				arrDeviceItem = append(arrDeviceItem, device)
			}
		}

		deviceAlarmList, err := inverter.GetDeviceAlarmList(plantCode, beginTime, endTime)
		if err != nil {
			r.logger.Error(err)
			return err
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

	var deviceIDList []string
	var deviceIDListString []string
	for _, device := range arrDeviceItem {
		if len(deviceIDList) == 100 {
			tmp := strings.Join(deviceIDList, ",")
			deviceIDListString = append(deviceIDListString, tmp)
			deviceIDList = []string{}
		}

		if device.ID != 0 {
			deviceIDList = append(deviceIDList, strconv.Itoa(device.ID))
		}
	}
	deviceIDListString = append(deviceIDListString, strings.Join(deviceIDList, ","))

	for _, deviceID := range deviceIDListString {
		realtimeDeviceData, err := inverter.GetDeviceRealtimeData(deviceID, "1")
		if err != nil {
			r.logger.Error(err)
			return err
		}

		for _, item := range realtimeDeviceData {
			if item.ID != 0 {
				mapInverterIDToRealtimeData[item.ID] = item
			}
		}
	}

	for _, plant := range plantList {
		for _, device := range mapPlantCodeToDevice[plant.Code] {
			if device.TypeID == 1 {
				realtimeDevice := mapInverterIDToRealtimeData[device.ID].DataItemMap

				if realtimeDevice.Status == 0 {
					shutdownTime := strconv.Itoa(int(endTime))
					if mapInverterIDToRealtimeData[device.ID].DataItemMap.InverterShutdown != nil {
						inverterShutdown, ok := realtimeDevice.InverterShutdown.(float64)
						if ok {
							shutdownTime = strconv.Itoa(int(inverterShutdown))
						}
					}

					key := fmt.Sprintf("Huawei,%s,%s,%s,%s", plant.Code, device.SerialNumber, device.Name, "Disconnect")
					val := fmt.Sprintf("%s,%s,%s", plant.Name, "Disconnect", shutdownTime)
					err := r.rdb.Set(ctx, key, val, 0).Err()
					if err != nil {
						r.logger.Error(err)
						return err
					}

					err = r.snmp.SendAlarmTrap(
						plant.Name,
						fmt.Sprintf("HUW-%s", "Disconnect"),
						fmt.Sprintf("Huawei,%s,%s", device.Name, "Disconnect"),
						"5",
						shutdownTime,
					)

					if err != nil {
						r.logger.Error(err)
						return err
					}

					continue
				}
			}

			if len(mapDeviceSNToAlarm[device.SerialNumber]) > 0 {
				for _, alarm := range mapDeviceSNToAlarm[device.SerialNumber] {
					alarmTime := strconv.Itoa(int(alarm.RaiseTime))
					key := fmt.Sprintf("Huawei,%s,%s,%s,%s", plant.Code, device.SerialNumber, device.Name, alarm.AlarmName)
					val := fmt.Sprintf("%s,%s,%s", plant.Name, alarm.AlarmCause, alarmTime)
					err := r.rdb.Set(ctx, key, val, 0).Err()
					if err != nil {
						r.logger.Error(err)
						return err
					}

					err = r.snmp.SendAlarmTrap(
						plant.Name,
						strings.ReplaceAll(
							fmt.Sprintf("HUW-%s", alarm.AlarmName),
							" ",
							"-",
						),
						fmt.Sprintf("Huawei,%s,%s", device.Name, alarm.AlarmCause),
						"5",
						alarm.AlarmName,
					)

					if err != nil {
						r.logger.Error(err)
						return err
					}
				}

				continue
			}

			var keys []string
			var cursor uint64
			for {
				var scanKeys []string
				scanKeys, cursor, err = r.rdb.Scan(ctx, cursor, fmt.Sprintf("Huawei,%s,%s,%s,*", plant.Code, device.SerialNumber, device.Name), 10).Result()
				if err != nil {
					r.logger.Error(err)
					return err
				}

				keys = append(keys, scanKeys...)
				if cursor == 0 {
					break
				}
			}

			for _, key := range keys {
				val, err := r.rdb.Get(ctx, key).Result()
				if err != nil {
					if err != redis.Nil {
						r.logger.Error(err)
						return err
					}

					continue
				}

				if val != "" {
					splitKey := strings.Split(key, ",")
					splitVal := strings.Split(val, ",")

					err = r.snmp.SendAlarmTrap(
						splitVal[0],
						strings.ReplaceAll(fmt.Sprintf("HUW-%s", splitKey[4]), " ", "-"),
						fmt.Sprintf("Huawei,%s,%s", device.Name, splitVal[1]),
						"0",
						splitVal[2],
					)

					if err != nil {
						return err
					}

					err = r.rdb.Del(ctx, key).Err()
					if err != nil {
						r.logger.Error(err)
						return err
					}
				}
			}
		}
	}

	return nil
}
