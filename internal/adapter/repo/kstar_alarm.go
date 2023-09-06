package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/kstar"
)

type kstarAlarmRepo struct {
	rdb    *redis.Client
	snmp   port.SNMPRepoPort
	logger logger.Logger
}

func NewKStarAlarmRepo(rdb *redis.Client, snmp port.SNMPRepoPort, logger logger.Logger) *kstarAlarmRepo {
	return &kstarAlarmRepo{rdb: rdb, snmp: snmp, logger: logger}
}

func (r kstarAlarmRepo) Alarm(username, password string) error {
	ctx := context.Background()
	inverter := kstar.NewKStarInverter(username, password, nil)

	deviceList, err := inverter.GetDeviceList()
	if err != nil {
		r.logger.Error((err))
		return err
	}

	for _, device := range deviceList {
		saveTime := device.SaveTime

		realtimeDeviceData, err := inverter.GetRealtimeDeviceData(device.ID)
		if err != nil {
			r.logger.Error(err)
			return err
		}

		if realtimeDeviceData != nil {
			saveTime = realtimeDeviceData.SaveTime
		}

		switch device.Status {
		case 0:
			key := fmt.Sprintf("Kstar,%s,%s,%s,%s", device.PlantID, device.ID, device.Name, "Kstar-Disconnect")
			val := fmt.Sprintf("%s,%s", device.PlantName, saveTime)
			err := r.rdb.Set(ctx, key, val, 0).Err()
			if err != nil {
				r.logger.Error(err)
				return err
			}

			err = r.snmp.SendAlarmTrap(
				device.PlantName,
				"Kstar-Disconnect",
				fmt.Sprintf("Kstar,%s,%s,%s", device.PlantID, device.ID, device.Name),
				"5",
				saveTime,
			)
			if err != nil {
				r.logger.Error(err)
				return err
			}
		case 1:
			realtimeAlarm, err := inverter.GetRealtimeAlarmListOfDevice(device.ID)
			if err != nil {
				r.logger.Error(err)
				return err
			}

			if len(realtimeAlarm) > 0 {
				err = r.snmp.SendAlarmTrap(
					device.PlantName,
					"Kstar-Disconnect",
					fmt.Sprintf("Kstar,%s,%s,%s", device.PlantID, device.ID, device.Name),
					"0",
					saveTime,
				)
				if err != nil {
					r.logger.Error(err)
					return err
				}

				err = r.rdb.Del(ctx, fmt.Sprintf("Kstar,%s,%s,%s,%s", device.PlantID, device.ID, device.Name, "Kstar-Disconnect")).Err()
				if err != nil {
					r.logger.Error(err)
					return err
				}

				for _, alarm := range realtimeAlarm {
					alarmTime := alarm.SaveTime
					alarmMessage := strings.ReplaceAll(alarm.Message, " ", "-")
					key := fmt.Sprintf("Kstar,%s,%s,%s,%s", device.PlantID, device.ID, device.Name, alarmMessage)
					val := fmt.Sprintf("%s,%s", device.PlantName, alarmTime)
					err := r.rdb.Set(ctx, key, val, 0).Err()
					if err != nil {
						r.logger.Error(err)
						return err
					}

					err = r.snmp.SendAlarmTrap(
						device.PlantName,
						alarmMessage,
						fmt.Sprintf("Kstar,%s,%s,%s", device.PlantID, device.ID, device.Name),
						"5",
						alarmTime,
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
				scanKeys, cursor, err = r.rdb.Scan(ctx, cursor, fmt.Sprintf("Kstar,%s,%s,%s,*", device.PlantID, device.ID, device.Name), 10).Result()
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
					r.logger.Error(err)
					return err
				}

				if !helper.EmptyString(val) {
					splitKey := strings.Split(key, ",")
					splitVal := strings.Split(val, ",")

					err = r.snmp.SendAlarmTrap(
						splitVal[0],
						strings.ReplaceAll(splitKey[4], " ", "-"),
						fmt.Sprintf("Kstar,%s,%s,%s", device.PlantID, device.ID, device.Name),
						"0",
						splitVal[1],
					)
					if err != nil {
						r.logger.Error(err)
						return err
					}

					err = r.rdb.Del(ctx, key).Err()
					if err != nil {
						r.logger.Error(err)
						return err
					}
				}
			}
		case 2:
			realtimeAlarm, err := inverter.GetRealtimeAlarmListOfDevice(device.ID)
			if err != nil {
				r.logger.Error(err)
				return err
			}

			if len(realtimeAlarm) > 0 {
				for _, alarm := range realtimeAlarm {
					alarmTime := alarm.SaveTime
					alarmMessage := strings.ReplaceAll(alarm.Message, " ", "-")

					key := fmt.Sprintf("Kstar,%s,%s,%s,%s", device.PlantID, device.ID, device.Name, alarmMessage)
					val := fmt.Sprintf("%s,%s", device.PlantName, alarmTime)
					err := r.rdb.Set(ctx, key, val, 0).Err()
					if err != nil {
						r.logger.Error(err)
						return err
					}

					err = r.snmp.SendAlarmTrap(
						device.PlantName,
						alarmMessage,
						fmt.Sprintf("Kstar,%s,%s,%s", device.PlantID, device.ID, device.Name),
						"5",
						alarmTime,
					)
					if err != nil {
						r.logger.Error(err)
						return err
					}
				}
			}
		default:
		}
	}

	return nil
}
