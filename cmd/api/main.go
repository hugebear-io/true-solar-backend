package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/api"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar-backend/pkg/middleware"
)

func main() {
	apiConfig := config.Config.API
	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/solar-api.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(apiConfig.LogLevel),
		SkipCaller:  1,
	})
	defer l.Close()

	// initialized database
	infra.InitDatabase(l)
	defer infra.SqlDB.Close()

	// api application
	app := gin.New()
	app.Use(middleware.CORS())
	router := app.Group("/api")

	// bind api
	api.BindHealthCheckAPI(router)
	api.BindAccessLogAPI(router, infra.SqlDB, l)
	api.BindAlarmConfigAPI(router, infra.SqlDB, l)
	api.BindDataCollectorConfigAPI(router, infra.SqlDB, l)
	api.BindInstalledCapacityConfigAPI(router, infra.SqlDB, l)
	api.BindLoginAPI(router, apiConfig.SecretKey, infra.SqlDB, l)
	api.BindPerformanceAlarmConfigAPI(router, infra.SqlDB, l)
	api.BindRedisConfigAPI(router, infra.SqlDB, l)
	api.BindSNMPConfigAPI(router, infra.SqlDB, l)
	api.BindUserAPI(router, infra.SqlDB)

	// launch
	addr := "0.0.0.0:3000"
	if apiConfig.Host != "" && apiConfig.Port != "" {
		addr = fmt.Sprintf("%v:%v", apiConfig.Host, apiConfig.Port)
	}

	l.Infof("server running on %v", addr)
	if err := app.Run(addr); err != nil {
		l.Fatal(err)
	}
}
