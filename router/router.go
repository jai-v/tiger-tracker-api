package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"tiger-tracker-api/configuration"
	"tiger-tracker-api/configuration/models"
	"tiger-tracker-api/constants"
	"tiger-tracker-api/controller"
	"tiger-tracker-api/docs"
	"tiger-tracker-api/logging"
)

func InitDB(poolConfig models.DbConnectionPool) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv(constants.DB_USER),
		os.Getenv(constants.DB_PASSWORD),
		os.Getenv(constants.DB_HOST),
		os.Getenv(constants.DB_PORT),
		os.Getenv(constants.DB_NAME))

	db, err := sqlx.Open("mysql", dsn)
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

	r := gin.New()
	appController := controller.NewAppController()

	routerGroup := r.Group("/api")
	{
		docs.SwaggerInfo.Title = "Tiger Tracker API"
		docs.SwaggerInfo.Description = "Tiger Tracker API Server for tracking the population of tigers in the wild"
		docs.SwaggerInfo.BasePath = "/api/tiger-tracker/v1"
		if configData.Environment != "prod" {
			routerGroup.GET("/tiger-tracker/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		routerGroup.GET("/tiger-tracker/health", appController.HealthCheck)
	}

	return r
}
