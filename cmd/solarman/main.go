package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/hugebear-io/true-solar-backend/api"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/job"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar-backend/pkg/middleware"
)

func init() {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	time.Local = loc
}

func init() {
	config.InitConfig()
}

func main() {
	apiConfig := config.Config.Solarman
	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/solarman-instance.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(apiConfig.LogLevel),
		SkipCaller:  1,
	})

	// initialized database
	infra.InitDatabase(l)

	cron := gocron.NewScheduler(time.UTC)

	collector, alarm := job.NewSolarmanJobHandler()
	cron.Cron("1 1-6 * * *").Do(collector.RunJob)
	cron.Cron("*/15 7-18 * * *").Do(alarm.RunJob)
	cron.StartAsync()

	// api application
	app := gin.New()
	app.Use(middleware.CORS())
	app.Use(gin.Recovery())
	router := app.Group("/api")

	// // bind api
	api.BindHealthCheckAPI(router)
	// api.BindSolarmanAPI(router)

	// // launch
	addr := "0.0.0.0:3002"
	if apiConfig.Host != "" && apiConfig.Port != "" {
		addr = fmt.Sprintf("%v:%v", apiConfig.Host, apiConfig.Port)
	}

	l.Infof("server running on %v", addr)
	if err := app.Run(addr); err != nil {
		l.Fatal(err)
	}
}
