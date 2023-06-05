package main

import (
	"fmt"
	"ledis/config"
	"ledis/lib/logger"
	"ledis/tcp"
	"os"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6380,
}

func fileExists(filename string) bool {
	stat, err := os.Stat(filename)
	return err == nil && stat.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "ledis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServerWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		}, tcp.NewHandler())
	if err != nil {
		logger.Error(err)
	}
}
