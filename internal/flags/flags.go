package flags

import (
	"flag"
	"os"

	"go.uber.org/zap"
)

var (
	FlagRunAddr     string
	FlagDatabaseURI string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagDatabaseURI, "d", ":9876", "address and port to run database")

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

	zap.L().Info(
		"Server settings",
		zap.String("Running server address: ", FlagRunAddr),
		zap.String("Running database address: ", FlagDatabaseURI),
	)
}
