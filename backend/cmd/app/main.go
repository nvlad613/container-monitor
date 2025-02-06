package main

import (
	"backend/config"
	"github.com/samber/lo"
)

func main() {
	conf := lo.Must(config.Load())

	logger := lo.Must(conf.Logger.Build()).Sugar()
	defer logger.Sync()

	logger.Infow("Service started!")
}
