package logging

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
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
		//logger.SetFormatter(&prefiex.TextFormatter{
		//
		//	ForceFormatting:  false,
		//	DisableTimestamp: false,
		//	DisableUppercase: false,
		//	FullTimestamp:    false,
		//	TimestampFormat:  "",
		//	DisableSorting:   false,
		//	QuoteEmptyFields: false,
		//	QuoteCharacter:   "",
		//	SpacePadding:     0,
		//	Once:             sync.Once{},
		//})
	}
}
