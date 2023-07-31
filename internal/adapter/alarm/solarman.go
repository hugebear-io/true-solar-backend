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
	logger logger.Logger,
) *solarmanAlarm {
	const brand = "INVT-Ipanda"
	return &solarmanAlarm{
		brand:        brand,
		rdb:          rdb,
		snmp:         snmp,
		usernameList: usernameList,
		password:     password,
		appID:        appID,
		appSecret:    appSecret,
		logger:       logger,
	}
}

func (r solarmanAlarm) Run() error {
	ctx := context.Background()
	now := time.Now()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location()) // 06:00 AM

	for _, username := range r.usernameList {
		inverter := solarman.NewSolarmanInverter(username, r.password, r.appID, r.appSecret, nil)

		basicTokenResp, err := inverter.GetBasicToken()
		if err != nil {
			return err
		}

		if helper.EmptyString(basicTokenResp.AccessToken) {
			return errors.New("access token is empty")
		}

		userInfoResp, err := inverter.GetUserInfo()
		if err != nil {
			return err
		}

		for _, company := range userInfoResp.OrgInfoList {
			businessTokenResp, err := inverter.GetBusinessToken(company.CompanyID)
			if err != nil {
				return err
			}

			if helper.EmptyString(businessTokenResp.AccessToken) {
				return errors.New("access token is empty")
			}

			token := businessTokenResp.AccessToken
			plantList, err := inverter.GetPlantList(token)
			for _, plant := range plantList {
				plantID := plant.ID
				plantName := plant.Name

				deviceList, err := inverter.GetPlantDeviceList(token, plantID)
				if err != nil {
					return err
				}

				for _, device := range deviceList {
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
							return err
						}

						name := fmt.Sprintf("%s-%s", plantName, deviceSN)
						alert := strings.ReplaceAll(fmt.Sprintf("%s-%s", deviceType, "Disconnect"), " ", "-")
						desc := fmt.Sprintf("%s,%d,%s,%d", r.brand, plantID, deviceSN, deviceID)
						severity := "5"
						err = r.snmp.SendAlarmTrap(name, alert, desc, severity, deviceCollectionTimeStr)
						if err != nil {
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
									return err
								}
							}

							err = r.rdb.Del(ctx, key).Err()
							if err != nil {
								return err
							}
						}
					case 2:
						alertList, err := inverter.GetDeviceAlertList(token, deviceSN, beginTime.Unix(), now.Unix())
						if err != nil {
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
									return err
								}

								name := fmt.Sprintf("%s-%s", plantName, deviceSN)
								alert := strings.ReplaceAll(fmt.Sprintf("%s-%s", deviceType, alertName), " ", "-")
								desc := fmt.Sprintf("%s,%d,%s,%d", r.brand, plantID, deviceSN, deviceID)
								severity := "5"

								err = r.snmp.SendAlarmTrap(name, alert, desc, severity, alertTimeStr)
								if err != nil {
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
	return nil
}
