package router

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"tiger-tracker-api/configuration"
	"tiger-tracker-api/configuration/models"
	"tiger-tracker-api/constants"
	"tiger-tracker-api/controller"
	"tiger-tracker-api/logging"
)

func InitDB(poolConfig models.DbConnectionPool) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv(constants.DB_USER),
		os.Getenv(constants.DB_PASSWORD),
		os.Getenv(constants.DB_HOST),
		os.Getenv(constants.DB_PORT),
		os.Getenv(constants.DB_NAME))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(poolConfig.MaxIdleConnections)
	db.SetMaxOpenConns(poolConfig.MaxOpenConnections)
	return db, nil
}

func Init(configData *configuration.ConfigData) *gin.Engine {
	logger := logging.GetLogger().WithField("Package", "router").WithField("Method", "Init")
	db, dbErr := InitDB(configData.DbConnectionPool)
	if dbErr != nil {
		logger.Fatalf("could not connect to db, error:%v", dbErr)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		logger.Fatalf("could not ping db, error:%v", pingErr)
	}

	router := gin.New()
	health := controller.NewHealthCheckController()
	router.GET("/tiger-tracker/api/health", health.Status)
	return router
}
