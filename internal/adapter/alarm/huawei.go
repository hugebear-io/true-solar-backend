package alarm

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

type huaweiAlarm struct {
	rdb          *redis.Client
	snmp         port.SNMPRepoPort
	usernameList []string
	password     string
	logger       logger.Logger
}

func NewHuaweiAlarm(
	rdb *redis.Client,
	snmp port.SNMPRepoPort,
	usernameList []string,
	password string,
	logger logger.Logger,
) *huaweiAlarm {
	return &huaweiAlarm{
		rdb:          rdb,
		snmp:         snmp,
		usernameList: usernameList,
		password:     password,
		logger:       logger,
	}
}

func (r huaweiAlarm) Run() error {
	now := time.Now()
	ctx := context.Background()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location()).UnixNano() / 1e6
	endTime := now.UnixNano() / 1e6

	for _, username := range r.usernameList {
		inverter := huawei.NewHuaweiInverter(username, r.password, nil)

		plantList, err := inverter.GetPlantList()
		if err != nil {
			r.logger.Error(err)
			return err
		}

		var plantCodeList []string
		var plantCodeListString []string
		for _, item := range plantList {
			if len(plantCodeList) == 100 {
				tmp := strings.Join(plantCodeList, ",")
				plantCodeListString = append(plantCodeListString, tmp)
				plantCodeList = []string{}
			}

			if !helper.EmptyString(item.Code) {
				plantCodeList = append(plantCodeList, item.Code)
			}
		}
		plantCodeListString = append(plantCodeListString, strings.Join(plantCodeList, ","))

		var deviceList []huawei.DeviceItem
		mapPlantCodeToDevice := make(map[string][]huawei.DeviceItem)
		mapDeviceSNToAlarm := make(map[string][]huawei.DeviceAlarmItem)
		mapInverterIDToRealtimeData := make(map[int]huawei.RealtimeDeviceData)

		for _, plantCode := range plantCodeListString {
			deviceListResp, err := inverter.GetDeviceList(plantCode)
			if err != nil {
				r.logger.Error(err)
				return err
			}

			for _, item := range deviceListResp {
				if !helper.EmptyString(item.PlantCode) {
					mapPlantCodeToDevice[item.PlantCode] = append(mapPlantCodeToDevice[item.PlantCode], item)
				}

				if item.TypeID == 1 {
					deviceList = append(deviceList, item)
				}
			}

			deviceAlarmListResp, err := inverter.GetDeviceAlarmList(plantCode, beginTime, endTime)
			if err != nil {
				r.logger.Error(err)
				return err
			}

			for _, alarm := range deviceAlarmListResp {
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

		for _, deviceID := range deviceIDListString {
			realtimeDeviceResp, err := inverter.GetDeviceRealtimeData(deviceID, "1")
			if err != nil {
				r.logger.Error(err)
				return err
			}

			for _, item := range realtimeDeviceResp {
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
							inverterShutdown, ok := (realtimeDevice.InverterShutdown).(float64)
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

						name := plant.Name
						alert := fmt.Sprintf("HUW-%s", "Disconnect")
						description := fmt.Sprintf("Huawei,%s,%s", device.Name, "Disconnect")
						severity := "5"

						if err := r.snmp.SendAlarmTrap(
							name,
							alert,
							description,
							severity,
							shutdownTime,
						); err != nil {
							r.logger.Error(err)
							return err
						}

						continue
					}
				}

				if len(mapDeviceSNToAlarm[device.SerialNumber]) > 0 {
					for _, alarm := range mapDeviceSNToAlarm[device.SerialNumber] {
						alarmTime := strconv.Itoa(int(alarm.RaiseTime))

						key := fmt.Sprintf(
							"Huawei,%s,%s,%s,%s",
							plant.Code,
							device.SerialNumber,
							device.Name,
							alarm.AlarmName,
						)

						val := fmt.Sprintf(
							"%s,%s,%s",
							plant.Name,
							alarm.AlarmCause,
							alarmTime,
						)

						err := r.rdb.Set(ctx, key, val, 0).Err()
						if err != nil {
							r.logger.Error(err)
							return err
						}

						name := plant.Name
						alert := strings.ReplaceAll(fmt.Sprintf("HUW-%s", alarm.AlarmName), " ", "-")
						description := fmt.Sprintf("Huawei,%s,%s", device.Name, alarm.AlarmCause)
						severity := "5"

						if err := r.snmp.SendAlarmTrap(name, alert, description, severity, alarmTime); err != nil {
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
					key := fmt.Sprintf("Huawei,%s,%s,%s,*", plant.Code, device.SerialNumber, device.Name)
					scanKeys, cursor, err = r.rdb.Scan(ctx, cursor, key, 10).Result()
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

					if !helper.EmptyString(val) {
						splitKey := strings.Split(key, ",")
						splitVal := strings.Split(val, ",")

						name := splitVal[0]
						alert := strings.ReplaceAll(fmt.Sprintf("HUW-%s", splitKey[4]), " ", "-")
						description := fmt.Sprintf("Huawei,%s,%s", device.Name, splitVal[1])
						severity := "0"
						alarmTime := splitVal[2]

						if err := r.snmp.SendAlarmTrap(
							name,
							alert,
							description,
							severity,
							alarmTime,
						); err != nil {
							r.logger.Error(err)
							return err
						}

						if err := r.rdb.Del(ctx, key).Err(); err != nil {
							r.logger.Error(err)
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
