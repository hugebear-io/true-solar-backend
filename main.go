package main

import (
	"fmt"
	"log"

	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

func init() {
	config.InitConfig()
}

func main() {
	l := logger.NewLoggerMock()
	elastic := infra.NewElasticSearch(l)
	repo := repo.NewElasticSearchRepo(elastic, "solarcell")
	data, err := repo.QueryPerformanceLow(30, 60, 24, 0.3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*data[0])
}
