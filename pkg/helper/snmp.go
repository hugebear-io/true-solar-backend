package helper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
)

func BuildSNMPPerformanceAlarmPayload(alarmType int, alarmConfig port.PerformanceAlarmConfig, capacityConfig port.InstalledCapacityConfig, data map[string]interface{}) (string, string, string, string, error) {
	if alarmType != constant.PERFORMANCE_ALARM_TYPE_PERFORMANCE_LOW && alarmType != constant.PERFORMANCE_ALARM_TYPE_SUM_PERFORMANCE_LOW {
		return "", "", "", "", errors.New("alarm type must be 1 (PERFORMANCE_LOW) or 2 (SUM_PERFORMANCE_LOW)")
	}

	var capacity float64
	if cap, ok := data["installedCapacity"].(float64); ok {
		capacity = cap
	}

	var plantItem port.PlantItem
	if err := Recast(data["plantItem"], &plantItem); err != nil {
		return "", "", "", "", err
	}

	var period string
	if p, ok := data["period"].(string); ok {
		period = p
	}

	var vendorName string
	switch strings.ToLower(plantItem.VendorType) {
	case constant.VENDOR_TYPE_GROWATT:
		vendorName = "Growatt"
	case constant.VENDOR_TYPE_HUAWEI:
		vendorName = "HUA"
	case constant.VENDOR_TYPE_KSTAR:
		vendorName = "Kstar"
	case constant.VENDOR_TYPE_INVT:
		vendorName = "INVT-Ipanda"
	case constant.VENDOR_TYPE_SOLARMAN: // Todo: Remove after change SOLARMAN to INVT in elasticsearch
		vendorName = "INVT-Ipanda"
	default:
		// no-op
	}

	if vendorName == "" {
		return "", "", "", "", fmt.Errorf("vendor type (%s) not supported", plantItem.VendorType)
	}

	plantName := plantItem.Name
	alarmName := fmt.Sprintf("SolarCell-%s", strings.ReplaceAll(alarmConfig.Name, " ", ""))
	alarmNameInDescription := AddSpace(alarmConfig.Name)
	severity := "5"
	duration := alarmConfig.Duration
	hitDay := alarmConfig.HitDay
	multipliedCapacity := capacity * capacityConfig.EfficiencyFactor * float64(capacityConfig.FocusHour)

	// PerformanceLow
	if alarmType == constant.PERFORMANCE_ALARM_TYPE_PERFORMANCE_LOW {
		severity := "3"
		payload := fmt.Sprintf("%s, %s, Less than or equal %.2f%%, Expected Daily Production:%.2f KWH, Actual Production less than:%.2f KWH, Duration:%d days, Period:%s",
			vendorName, alarmNameInDescription, alarmConfig.Percentage, multipliedCapacity, multipliedCapacity*(alarmConfig.Percentage/100.0), hitDay, period)
		return plantName, alarmName, payload, severity, nil
	}

	// SumPerformanceLow
	var totalProduction float64
	if x, ok := data["totalProduction"].(float64); ok {
		totalProduction = x
	}

	payload := fmt.Sprintf("%s, %s, Less than or equal %.2f%%, Expected Production:%.2f KWH, Actual Production:%.2f KWH (less than %.2f KWH), Duration:%d days, Period:%s",
		vendorName, alarmNameInDescription, alarmConfig.Percentage, multipliedCapacity*float64(duration), totalProduction, (multipliedCapacity*float64(duration))*(alarmConfig.Percentage/100.0), duration, period)
	return plantName, alarmName, payload, severity, nil
}
