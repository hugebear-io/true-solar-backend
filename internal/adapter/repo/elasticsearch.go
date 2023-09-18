package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/olivere/elastic/v7"
)

type elasticSearchRepo struct {
	index  string
	client *elastic.Client
}

func NewElasticSearchRepo(client *elastic.Client, index string) port.ElasticSearchRepoPort {
	es := elasticSearchRepo{
		client: client,
		index:  index,
	}

	return &es
}

func (es elasticSearchRepo) createIndexIfNotExist(index string) error {
	ctx := context.Background()
	if exists, err := es.client.IndexExists(index).Do(ctx); err == nil {
		if !exists {
			result, err := es.client.CreateIndex(index).Do(ctx)
			if err != nil {
				return err
			}

			if !result.Acknowledged {
				err := errors.New("elastic search did not acknowledge")
				return err
			}
		}
	}
	return nil
}

// BulkIndex indexes documents
func (es elasticSearchRepo) BulkIndex(docs []interface{}) error {
	now := time.Now().UTC()
	index := fmt.Sprintf("%s-%s", es.index, now.Format("2006.01.02"))

	err := es.createIndexIfNotExist(index)
	if err != nil {
		return err
	}

	bulkRequest := es.client.Bulk()

	for _, doc := range docs {
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(index).Doc(doc))
	}

	_, err = bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (es elasticSearchRepo) UpsertSiteStation(docs []port.SiteItem) error {
	index := "site-station"

	err := es.createIndexIfNotExist(index)
	if err != nil {
		return err
	}

	bulkRequest := es.client.Bulk()

	for _, doc := range docs {
		bulkRequest.Add(elastic.NewBulkUpdateRequest().Index(index).Id(doc.SiteID).Doc(doc).DocAsUpsert(true))
	}

	_, err = bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (es elasticSearchRepo) DeleteSiteStationDocument() error {
	index := "site-station"
	_, err := es.client.DeleteByQuery(index).Query(&elastic.MatchAllQuery{}).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// SearchIndex returns search index service
func (es elasticSearchRepo) SearchIndex() *elastic.SearchService {
	return es.client.Search(fmt.Sprintf("%s*", es.index))
}

// QueryPerformanceLow returns plant item within performance low conditions
func (es elasticSearchRepo) QueryPerformanceLow(duration int, EfficiencyFactor float64, FocusHour int, thresholdPercentage float64) ([]*elastic.AggregationBucketCompositeItem, error) {
	ctx := context.Background()
	items := make([]*elastic.AggregationBucketCompositeItem, 0)

	compositeAggregation := elastic.NewCompositeAggregation().
		Size(10000).
		Sources(elastic.NewCompositeAggregationDateHistogramValuesSource("date").Field("@timestamp").CalendarInterval("1d").Format("yyyy-MM-dd"),
			elastic.NewCompositeAggregationTermsValuesSource("vendor_type").Field("vendor_type.keyword"),
			elastic.NewCompositeAggregationTermsValuesSource("id").Field("id.keyword")).
		SubAggregation("max_daily", elastic.NewMaxAggregation().Field("daily_production")).
		SubAggregation("avg_capacity", elastic.NewAvgAggregation().Field("installed_capacity")).
		SubAggregation("threshold_percentage", elastic.NewBucketScriptAggregation().
			BucketsPathsMap(map[string]string{"capacity": "avg_capacity"}).
			Script(elastic.NewScript("params.capacity * params.efficiency_factor * params.focus_hour * params.threshold_percentage").
				Params(map[string]interface{}{
					"efficiency_factor":    EfficiencyFactor,
					"focus_hour":           FocusHour,
					"threshold_percentage": thresholdPercentage,
				}))).
		SubAggregation("under_threshold", elastic.NewBucketSelectorAggregation().
			BucketsPathsMap(map[string]string{"threshold": "threshold_percentage", "daily": "max_daily"}).
			Script(elastic.NewScript("params.daily <= params.threshold"))).
		SubAggregation("hits", elastic.NewTopHitsAggregation().
			Size(1).
			FetchSourceContext(
				elastic.NewFetchSourceContext(true).Include(
					"id", "name", "vendor_type", "node_type", "ac_phase", "plant_status",
					"area", "site_id", "site_city_code", "site_city_name", "installed_capacity",
				)))

	searchQuery := es.SearchIndex().
		Size(0).
		Query(elastic.NewBoolQuery().Must(
			elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
			elastic.NewRangeQuery("@timestamp").Gte(fmt.Sprintf("now-%dd/d", duration)).Lte("now-1d/d"),
		)).
		Aggregation("performance_alarm", compositeAggregation)

	firstResult, err := searchQuery.Pretty(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	if firstResult.Aggregations == nil {
		return nil, errors.New("cannot get result aggregations")
	}

	performanceAlarm, found := firstResult.Aggregations.Composite("performance_alarm")
	if !found {
		return nil, errors.New("cannot get result composite performance alarm")
	}

	items = append(items, performanceAlarm.Buckets...)

	if len(performanceAlarm.AfterKey) > 0 {
		afterKey := performanceAlarm.AfterKey

		for {
			searchQuery = es.SearchIndex().
				Size(0).
				Query(elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
					elastic.NewRangeQuery("@timestamp").Gte(fmt.Sprintf("now-%dd/d", duration)).Lte("now-1d/d"),
				)).
				Aggregation("performance_alarm", compositeAggregation.AggregateAfter(afterKey))

			result, err := searchQuery.Pretty(true).Do(ctx)
			if err != nil {
				return nil, err
			}

			if firstResult.Aggregations == nil {
				return nil, errors.New("cannot get result aggregations")
			}

			performanceAlarm, found := result.Aggregations.Composite("performance_alarm")
			if !found {
				return nil, errors.New("cannot get result composite performance alarm")
			}

			items = append(items, performanceAlarm.Buckets...)

			if len(performanceAlarm.AfterKey) == 0 {
				break
			}

			afterKey = performanceAlarm.AfterKey
		}
	}

	return items, err
}

// QuerySumPerformanceLow returns plant item within sum performance low conditions
func (es elasticSearchRepo) QuerySumPerformanceLow(duration int) ([]*elastic.AggregationBucketCompositeItem, error) {
	ctx := context.Background()
	items := make([]*elastic.AggregationBucketCompositeItem, 0)

	compositeAggregation := elastic.NewCompositeAggregation().
		Size(10000).
		Sources(elastic.NewCompositeAggregationDateHistogramValuesSource("date").Field("@timestamp").CalendarInterval("1d").Format("yyyy-MM-dd"),
			elastic.NewCompositeAggregationTermsValuesSource("vendor_type").Field("vendor_type.keyword"),
			elastic.NewCompositeAggregationTermsValuesSource("id").Field("id.keyword")).
		SubAggregation("max_daily", elastic.NewMaxAggregation().Field("daily_production")).
		SubAggregation("avg_capacity", elastic.NewAvgAggregation().Field("installed_capacity")).
		SubAggregation("hits", elastic.NewTopHitsAggregation().
			Size(1).
			FetchSourceContext(
				elastic.NewFetchSourceContext(true).Include(
					"id", "name", "vendor_type", "node_type", "ac_phase", "plant_status",
					"area", "site_id", "site_city_code", "site_city_name", "installed_capacity",
				)))

	searchQuery := es.SearchIndex().
		Size(0).
		Query(elastic.NewBoolQuery().Must(
			elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
			elastic.NewRangeQuery("@timestamp").Gte(fmt.Sprintf("now-%dd/d", duration)).Lte("now-1d/d"),
		)).
		Aggregation("performance_alarm", compositeAggregation)

	firstResult, err := searchQuery.Pretty(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	if firstResult.Aggregations == nil {
		return nil, errors.New("cannot get result aggregations")
	}

	performanceAlarm, found := firstResult.Aggregations.Composite("performance_alarm")
	if !found {
		return nil, errors.New("cannot get result composite performance alarm")
	}

	items = append(items, performanceAlarm.Buckets...)

	if len(performanceAlarm.AfterKey) > 0 {
		afterKey := performanceAlarm.AfterKey

		for {
			searchQuery = es.SearchIndex().
				Size(0).
				Query(elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
					elastic.NewRangeQuery("@timestamp").Gte(fmt.Sprintf("now-%dd/d", duration)).Lte("now-1d/d"),
				)).
				Aggregation("performance_alarm", compositeAggregation.AggregateAfter(afterKey))

			result, err := searchQuery.Pretty(true).Do(ctx)
			if err != nil {
				return nil, err
			}

			if firstResult.Aggregations == nil {
				return nil, errors.New("cannot get result aggregations")
			}

			performanceAlarm, found := result.Aggregations.Composite("performance_alarm")
			if !found {
				return nil, errors.New("cannot get result composite performance alarm")
			}

			items = append(items, performanceAlarm.Buckets...)

			if len(performanceAlarm.AfterKey) == 0 {
				break
			}

			afterKey = performanceAlarm.AfterKey
		}
	}

	return items, err
}

func (es elasticSearchRepo) QueryPerformanceOK(duration int, EfficiencyFactor float64, FocusHour int, thresholdPercentage float64) ([]*elastic.AggregationBucketCompositeItem, error) {
	ctx := context.Background()
	items := make([]*elastic.AggregationBucketCompositeItem, 0)

	compositeAggregation := elastic.NewCompositeAggregation().
		Size(10000).
		Sources(elastic.NewCompositeAggregationDateHistogramValuesSource("date").Field("@timestamp").CalendarInterval("1d").Format("yyyy-MM-dd"),
			elastic.NewCompositeAggregationTermsValuesSource("vendor_type").Field("vendor_type.keyword"),
			elastic.NewCompositeAggregationTermsValuesSource("id").Field("id.keyword")).
		SubAggregation("max_daily", elastic.NewMaxAggregation().Field("daily_production")).
		SubAggregation("avg_capacity", elastic.NewAvgAggregation().Field("installed_capacity")).
		SubAggregation("threshold_percentage", elastic.NewBucketScriptAggregation().
			BucketsPathsMap(map[string]string{"capacity": "avg_capacity"}).
			Script(elastic.NewScript("params.capacity * params.efficiency_factor * params.focus_hour * params.threshold_percentage").
				Params(map[string]interface{}{
					"efficiency_factor":    EfficiencyFactor,
					"focus_hour":           FocusHour,
					"threshold_percentage": thresholdPercentage,
				}))).
		SubAggregation("above_threshold", elastic.NewBucketSelectorAggregation().
			BucketsPathsMap(map[string]string{"threshold": "threshold_percentage", "daily": "max_daily"}).
			Script(elastic.NewScript("params.daily >= params.threshold"))).
		SubAggregation("hits", elastic.NewTopHitsAggregation().
			Size(1).
			FetchSourceContext(
				elastic.NewFetchSourceContext(true).Include(
					"id", "name", "vendor_type", "node_type", "ac_phase", "plant_status",
					"area", "site_id", "site_city_code", "site_city_name", "installed_capacity",
				)))

	searchQuery := es.SearchIndex().
		Size(0).
		Query(elastic.NewBoolQuery().Must(
			elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
			elastic.NewRangeQuery("@timestamp").Gte(fmt.Sprintf("now-%dd/d", duration)).Lte("now-1d/d"),
		)).
		Aggregation("performance_alarm", compositeAggregation)

	firstResult, err := searchQuery.Pretty(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	if firstResult.Aggregations == nil {
		return nil, errors.New("cannot get result aggregations")
	}

	performanceAlarm, found := firstResult.Aggregations.Composite("performance_alarm")
	if !found {
		return nil, errors.New("cannot get result composite performance alarm")
	}

	items = append(items, performanceAlarm.Buckets...)

	if len(performanceAlarm.AfterKey) > 0 {
		afterKey := performanceAlarm.AfterKey

		for {
			searchQuery = es.SearchIndex().
				Size(0).
				Query(elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
					elastic.NewRangeQuery("@timestamp").Gte(fmt.Sprintf("now-%dd/d", duration)).Lte("now-1d/d"),
				)).
				Aggregation("performance_alarm", compositeAggregation.AggregateAfter(afterKey))

			result, err := searchQuery.Pretty(true).Do(ctx)
			if err != nil {
				return nil, err
			}

			if firstResult.Aggregations == nil {
				return nil, errors.New("cannot get result aggregations")
			}

			performanceAlarm, found := result.Aggregations.Composite("performance_alarm")
			if !found {
				return nil, errors.New("cannot get result composite performance alarm")
			}

			items = append(items, performanceAlarm.Buckets...)

			if len(performanceAlarm.AfterKey) == 0 {
				break
			}

			afterKey = performanceAlarm.AfterKey
		}
	}

	return items, err
}
