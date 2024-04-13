package main

import (
	"tiger-tracker-api/configuration"
	"tiger-tracker-api/constants"
	"tiger-tracker-api/logging"
	"tiger-tracker-api/router"
)

func main() {
	logger := logging.GetLogger().WithField("Package", "main").WithField("Method", "main")
	logger.Info("starting service")
	configData, loadConfigErr := configuration.NewConfigLoader().LoadConfig(constants.CONFIG_FILE_NAME)
	if loadConfigErr != nil {
		logger.Fatalf("could not load config, error:%v", loadConfigErr)
	}

	r := router.Init(configData)
	err := r.Run(":" + configData.ListenPort)
	if err != nil {
		logger.Fatalf("encountered fatal error: %v", err)
	}
}
