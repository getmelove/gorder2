package logging

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	prefiex "github.com/x-cray/logrus-prefixed-formatter"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
		CallerPrettyfier: nil,
		PrettyPrint:      false,
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		logger.SetFormatter(&prefiex.TextFormatter{
			ForceColors:      false,
			DisableColors:    false,
			ForceFormatting:  true,
			DisableTimestamp: false,
			DisableUppercase: false,
			FullTimestamp:    false,
			TimestampFormat:  "",
			DisableSorting:   false,
			QuoteEmptyFields: false,
			QuoteCharacter:   "",
			SpacePadding:     0,
			Once:             sync.Once{},
		})
	}
}
