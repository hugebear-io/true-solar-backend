package alarm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/solarman"
)

// bignode.invt.th@gmail.com
// 123456*
type solarmanAlarm struct {
	rdb          *redis.Client
	snmp         port.SNMPRepoPort
	brand        string
	usernameList []string
	password     string
	appID        string
	appSecret    string
	logger       logger.Logger
}

func NewSolarmanAlarm(
	rdb *redis.Client,
	snmp port.SNMPRepoPort,
	usernameList []string,
	password string,
	appID string,
	appSecret string,
) *solarmanAlarm {
	const brand = "INVT-Ipanda"

	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/solarman-alarm.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(logger.LOG_LEVEL_DEBUG),
		SkipCaller:  1,
	})

	return &solarmanAlarm{
		brand:        brand,
		rdb:          rdb,
		snmp:         snmp,
		usernameList: usernameList,
		password:     password,
		appID:        appID,
		logger:       l,
		appSecret:    appSecret,
	}
}

func (r solarmanAlarm) Run() error {
	defer r.logger.Close()
	ctx := context.Background()
	now := time.Now()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location()) // 06:00 AM

	totalUser := len(r.usernameList)
	for numUser, username := range r.usernameList {
		r.logger.Infof("User Count: %d/%d", numUser+1, totalUser)

		inverter := solarman.NewSolarmanInverter(username, r.password, r.appID, r.appSecret, nil)
		basicTokenResp, err := inverter.GetBasicToken()
		if err != nil {
			r.logger.Error(err)
			return err
		}

		if helper.EmptyString(basicTokenResp.AccessToken) {
			err := errors.New("access token is empty")
			r.logger.Error(err)
			return err
		}

		userInfoResp, err := inverter.GetUserInfo()
		if err != nil {
			r.logger.Error(err)
			return err
		}

		totalCompany := len(userInfoResp.OrgInfoList)
		for numCompany, company := range userInfoResp.OrgInfoList {
			r.logger.Infof("Company Count: %d/%d", numCompany+1, totalCompany)

			businessTokenResp, err := inverter.GetBusinessToken(company.CompanyID)
			if err != nil {
				r.logger.Error(err)
				return err
			}

			if helper.EmptyString(businessTokenResp.AccessToken) {
				err := errors.New("access token is empty")
				r.logger.Error(err)
				return err
			}

			token := businessTokenResp.AccessToken
			plantList, err := inverter.GetPlantList(token)
			if err != nil {
				r.logger.Error(err)
				return err
			}

			totalPlant := len(plantList)
			for numPlant, plant := range plantList {
				r.logger.Infof("Plant Count: %d/%d", numPlant+1, totalPlant)

				plantID := plant.ID
				plantName := plant.Name

				deviceList, err := inverter.GetPlantDeviceList(token, plantID)
				if err != nil {
					r.logger.Error(err)
					return err
				}

				totalDevice := len(deviceList)
				for numDevice, device := range deviceList {
					r.logger.Infof("Device Count: %d/%d", numDevice+1, totalDevice)

					deviceID := device.DeviceID
					deviceSN := device.DeviceSN
					deviceType := device.DeviceType
					deviceCollectionTime := device.CollectionTime
					deviceCollectionTimeStr := strconv.FormatInt(deviceCollectionTime, 10)

					switch device.ConnectStatus {
					case 0:
						key := fmt.Sprintf("%s,%d,%s,%s,%d,%s", r.brand, plantID, deviceType, deviceSN, deviceID, "Disconnect")
						val := fmt.Sprintf("%s,%s", plantName, deviceCollectionTimeStr)
						err := r.rdb.Set(ctx, key, val, 0).Err()
						if err != nil {
							r.logger.Error(err)
							return err
						}

						name := fmt.Sprintf("%s-%s", plantName, deviceSN)
						alert := strings.ReplaceAll(fmt.Sprintf("%s-%s", deviceType, "Disconnect"), " ", "-")
						desc := fmt.Sprintf("%s,%d,%s,%d", r.brand, plantID, deviceSN, deviceID)
						severity := "5"
						err = r.snmp.SendAlarmTrap(name, alert, desc, severity, deviceCollectionTimeStr)
						if err != nil {
							r.logger.Error(err)
							return err
						}
					case 1:
						var keys []string
						var cursor uint64
						for {
							var scanKeys []string
							key := fmt.Sprintf("%s,%d,%s,%s,%d,*", r.brand, plantID, deviceType, deviceSN, deviceID)
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
							if err == redis.Nil {
								continue
							}

							if err != nil {
								r.logger.Error(err)
								return err
							}

							if !helper.EmptyString(val) {
								splitKey := strings.Split(key, ",")
								splitVal := strings.Split(val, ",")

								name := fmt.Sprintf("%s-%s", plantName, deviceSN)
								alert := strings.ReplaceAll(fmt.Sprintf("%s-%s", deviceType, splitKey[5]), " ", "-")
								desc := fmt.Sprintf("%s,%d,%s,%d", r.brand, plantID, deviceSN, deviceID)
								severity := "0"

								err := r.snmp.SendAlarmTrap(name, alert, desc, severity, splitVal[1])
								if err != nil {
									r.logger.Error(err)
									return err
								}
							}

							err = r.rdb.Del(ctx, key).Err()
							if err != nil {
								r.logger.Error(err)
								return err
							}
						}
					case 2:
						alertList, err := inverter.GetDeviceAlertList(token, deviceSN, beginTime.Unix(), now.Unix())
						if err != nil {
							r.logger.Error(err)
							return err
						}

						for _, alert := range alertList {
							alertName := alert.AlertNameInPAAS
							alertTime := alert.AlertTime
							alertTimeStr := strconv.FormatInt(alertTime, 10)

							if !helper.EmptyString(alert.AlertNameInPAAS) && alert.AlertTime != 0 {
								key := fmt.Sprintf("%s,%d,%s,%s,%d,%s", r.brand, plantID, deviceType, deviceSN, deviceID, alertName)
								val := fmt.Sprintf("%s,%s", plantName, alertTimeStr)
								err := r.rdb.Set(ctx, key, val, 0).Err()
								if err != nil {
									r.logger.Error(err)
									return err
								}

								name := fmt.Sprintf("%s-%s", plantName, deviceSN)
								alert := strings.ReplaceAll(fmt.Sprintf("%s-%s", deviceType, alertName), " ", "-")
								desc := fmt.Sprintf("%s,%d,%s,%d", r.brand, plantID, deviceSN, deviceID)
								severity := "5"

								err = r.snmp.SendAlarmTrap(name, alert, desc, severity, alertTimeStr)
								if err != nil {
									r.logger.Error(err)
									return err
								}
							}
						}
					default:
					}
				}
			}
		}
	}

	r.logger.Info("Alarm Done")
	return nil
}
