package domain

import (
	"fmt"
	"strings"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
)

func BuildSNMPClearPerformanceAlarmPayload(alarmConfig port.PerformanceAlarmConfig, capacityConfig port.InstalledCapacityConfig, data map[string]interface{}) (string, string, string, string, error) {
	var capacity float64
	if cap, ok := data["installedCapacity"].(float64); ok {
		capacity = cap
	}

	var plantItem port.PlantItem
	if err := helper.Recast(data["plantItem"], &plantItem); err != nil {
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
	alarmNameInDescription := helper.AddSpace(alarmConfig.Name)
	severity := constant.CLEAR_SEVERITY
	hitDay := alarmConfig.HitDay
	multipliedCapacity := capacity * capacityConfig.EfficiencyFactor * float64(capacityConfig.FocusHour)

	payload := fmt.Sprintf("%s, %s, Greater than or equal %.2f%%, Expected Daily Production:%.2f KWH, Actual Production Greater than:%.2f KWH, Duration:%d days, Period:%s",
		vendorName, alarmNameInDescription, alarmConfig.Percentage, multipliedCapacity, multipliedCapacity*(alarmConfig.Percentage/100.0), hitDay, period)
	return plantName, alarmName, payload, severity, nil

}
