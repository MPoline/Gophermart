package main

import (
	"fmt"
	"os"

	"github.com/MPoline/Gophermart/internal/api"
	"github.com/MPoline/Gophermart/internal/database"
	"github.com/MPoline/Gophermart/internal/flags"
	"github.com/MPoline/Gophermart/internal/logging"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.InitLog()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing logger:", err)
	}
	defer logger.Sync()

	flags.ParseFlags()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	err = database.DBInit()
	if err != nil {
		logger.Warn("Error start database: ", zap.Error(err))
	}

	r := api.InitRouter()

	err = r.Run(flags.FlagRunAddr)
	if err != nil {
		logger.Warn("Error start server: ", zap.Error(err))
	}

}
