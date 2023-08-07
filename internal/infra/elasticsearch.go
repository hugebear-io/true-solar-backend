package infra

import (
	"crypto/tls"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/olivere/elastic/v7"
)

var httpsRegexp = regexp.MustCompile("^https")

func NewElasticSearch(logger logger.Logger) *elastic.Client {
	cfg := config.Config.ElasticSearch
	logger.Debugf("%v", cfg)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  true,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}

	scheme := "http"
	if httpsRegexp.FindString(cfg.Host) != "" {
		scheme = "https"
	}

	var err error
	elastic, err := elastic.NewClient(
		elastic.SetURL(cfg.Host),
		elastic.SetScheme(scheme),
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient),
		elastic.SetHealthcheckTimeout(time.Duration(300)*time.Second),
	)

	if err != nil {
		logger.Fatal(err)
		return nil
	}

	logger.Info("Initialized Elastic Search")
	return elastic
}
