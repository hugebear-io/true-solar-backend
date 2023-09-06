package infra

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

var SqlDB *sql.DB

func InitDatabase(logger logger.Logger) {
	cfg := config.Config.Database
	var err error
	SqlDB, err = sql.Open("sqlite3", cfg.DSN)
	if err != nil {
		logger.Panic(err)
	}

	logger.Info("Initial Database")
}
