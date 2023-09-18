package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/helper"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"go.openly.dev/pointy"
)

type clearPerformanceAlarmService struct {
	elastic                 port.ElasticSearchRepoPort
	snmp                    port.SNMPRepoPort
	performanceAlarmConfig  port.PerformanceAlarmConfig
	installedCapacityConfig port.InstalledCapacityConfig
	logger                  logger.Logger
}

func NewClearPerformanceAlarmService(elastic port.ElasticSearchRepoPort) *clearPerformanceAlarmService {
	return &clearPerformanceAlarmService{elastic: elastic}
}

func (s *clearPerformanceAlarmService) Run() {
	t := time.Now()
	s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): start performance low alarm at %s", t)

	if s.performanceAlarmConfig.HitDay == 0 {
		s.logger.Warn("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): cannot find hitday in performance alarm configuration or performance alarm's hitday must not be zero")
		return
	}

	if s.performanceAlarmConfig.Duration == 0 {
		s.logger.Warn("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): cannot find duration in performance alarm configuration or performance alarm's duration must not be zero")
		return
	}

	hitDay := s.performanceAlarmConfig.HitDay
	duration := s.performanceAlarmConfig.Duration
	percentage := s.performanceAlarmConfig.Percentage

	buckets, err := s.elastic.QueryPerformanceLow(duration, s.installedCapacityConfig.EfficiencyFactor, s.installedCapacityConfig.FocusHour, percentage/100.0)
	if err != nil {
		s.logger.Errorf("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): %s", err.Error())
	}

	period := fmt.Sprintf("%s - %s", t.UTC().AddDate(0, 0, -duration).Format("02Jan2006"), t.UTC().AddDate(0, 0, -1).Format("02Jan2006"))

	filteredBuckets := make(map[string]map[string]interface{}) // count, installedCapacity, *plantItem, period
	for _, bucketPtr := range buckets {
		if bucketPtr != nil {
			bucket := *bucketPtr

			if len(bucket.Key) == 0 {
				continue
			}

			var plantItem *port.PlantItem = nil
			var key string
			var installedCapacity float64

			if len(bucket.Key) > 0 {
				if vendorType, ok := bucket.Key["vendor_type"].(string); ok {
					if id, ok := bucket.Key["id"].(string); ok {
						key = fmt.Sprintf("%s,%s", vendorType, id)
					}
				}
			}

			if avgCapacity, ok := bucket.ValueCount("avg_capacity"); ok {
				installedCapacity = pointy.Float64Value(avgCapacity.Value, 0.0)
			}

			if topHits, found := bucket.Aggregations.TopHits("hits"); found {
				if topHits.Hits != nil {
					if len(topHits.Hits.Hits) == 1 {
						searchHitPtr := topHits.Hits.Hits[0]
						if searchHitPtr != nil {
							if err := json.Unmarshal(searchHitPtr.Source, &plantItem); err != nil {
								s.logger.Errorf("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): cannot unmarshal plant item document, (%s)", err)
								continue
							}
						}
					}
				}
			}

			if key != "" {
				if item, found := filteredBuckets[key]; found {
					if count, ok := item["count"].(int); ok {
						item["count"] = count + 1
					}
					filteredBuckets[key] = item
				} else {
					filteredBuckets[key] = map[string]interface{}{
						"count":             1,
						"installedCapacity": installedCapacity,
						"plantItem":         plantItem,
						"period":            period,
					}
				}
			}
		}
	}

	s.logger.Info("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): start performance low alarm stats...")
	s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): unfiltered %d plants", len(filteredBuckets))

	var alarmCount int
	var failedAlarmCount int
	if len(filteredBuckets) > 0 {
		bucketBatches := helper.ChunkBy(filteredBuckets, constant.PERFORMANCE_ALARM_SNMP_BATCH_SIZE)
		s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): %d batches (size of %d)", len(bucketBatches), constant.PERFORMANCE_ALARM_SNMP_BATCH_SIZE)

		var batchAlarmCount int
		var failedBatchAlarmCount int
		for i, batches := range bucketBatches {
			batchAlarmCount = 0
			failedBatchAlarmCount = 0

			for _, batch := range batches {
				for _, data := range batch {
					if count, ok := data["count"].(int); ok {
						if count >= hitDay {
							plantName, alarmName, alarmDescription, severity, err := domain.BuildSNMPClearPerformanceAlarmPayload(s.performanceAlarmConfig, s.installedCapacityConfig, data)
							if err != nil {
								s.logger.Errorf("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): %s", err.Error())
								continue
							}

							err = s.snmp.SendAlarmTrap(plantName, alarmName, alarmDescription, severity, t.UTC().Format(time.RFC3339Nano))
							if err != nil {
								failedAlarmCount++
								failedBatchAlarmCount++
								s.logger.Errorf("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): %s", err.Error())
								continue
							}

							alarmCount++
							batchAlarmCount++
						}
					}
				}
			}

			s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): batch %d - sending %d alarms (%d alarms failed sending to SNMP)", i+1, batchAlarmCount, failedBatchAlarmCount)
			s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): wait for %.2fs", constant.PERFORMANCE_ALARM_SNMP_BATCH_DELAY.Seconds())
			time.Sleep(constant.PERFORMANCE_ALARM_SNMP_BATCH_DELAY)
		}
	}

	s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): overall - sending %d alarms (%d alarms failed sending to SNMP)", alarmCount, failedAlarmCount)
	s.logger.Infof("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): polling finished in %s", time.Since(t).String())
	s.logger.Info("PerformanceAlarmScheduler.StartPerformanceLowAlarm(): end performance alarm low stats")

}
