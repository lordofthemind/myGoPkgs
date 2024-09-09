package mygopkgs

import (
	"os"

	"github.com/lordofthemind/myGoPkgs/logger"
)

func SetUpLoggerFile(logFileName string) (*os.File, error) {
	logfile, err := logger.LoggerSetUpLoggerFile(logFileName)
	return logfile, err
}
