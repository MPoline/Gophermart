package flags

import (
	"flag"
	"os"

	"go.uber.org/zap"
)

var (
	FlagRunAddr              string
	FlagDatabaseURI          string
	FlagAccuralSystemAddress string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagDatabaseURI, "d", ":9876", "address and port to run database")
	flag.StringVar(&FlagAccuralSystemAddress, "r", ":9000", "address and port to run accural")

	flag.Parse()

	if flag.NArg() > 0 {
		zap.L().Info("Error: unknown flag(s)")
		flag.Usage()
		return
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		zap.L().Info("ADDRESS: ", zap.String("envRunAddr", envRunAddr))
		FlagRunAddr = envRunAddr
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		zap.L().Info("DATABASE_URI: ", zap.String("envDatabaseURI", envDatabaseURI))
		FlagDatabaseURI = envDatabaseURI
	}

	if envAccuralSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccuralSystemAddress != "" {
		zap.L().Info("ACCRUAL_SYSTEM_ADDRESS: ", zap.String("envAccuralSystemAddress", envAccuralSystemAddress))
		FlagAccuralSystemAddress = envAccuralSystemAddress
	}

	zap.L().Info(
		"Server settings",
		zap.String("Running server address: ", FlagRunAddr),
		zap.String("Running database address: ", FlagDatabaseURI),
		zap.String("Running accural address: ", FlagAccuralSystemAddress),
	)
}
