package jokit

import (
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var appName string

type LoggerConfig struct {
	AppEnv    string
	AppName   string
	LogFolder string
	LogLevel  string
}

// InitLogger Initializes logger, writes log messages to files
// * If you want to log locally (e.g in development) set the env value to either (LOCAL | DEV)
// Log files are named from the APP_NAME and a time stamp attached to rotated files
// * if env is LOCAL, then logs will be output to stdErr
// * this is got from :: https://gitlab.betika.private/betikateam/gokit/-/blob/fff076dec2d179482352f1e513b22c7fc53a3a38/logger.go#L53
// * in case of any difference in opinion, feel free to raise and change
func InitLogger(loggerConfig LoggerConfig) error {
	appName = loggerConfig.AppName
	if strings.EqualFold(loggerConfig.AppEnv, "LOCAL") {
		log.SetOutput(os.Stderr)
		return nil
	}
	if strings.EqualFold(loggerConfig.AppEnv, "DEV") {
		pwd, err := os.Getwd()
		if err == nil {
			loggerConfig.LogFolder = fmt.Sprintf("%s/logs/", pwd)
		}
	}
	if strings.EqualFold(loggerConfig.LogFolder, "") {
		return fmt.Errorf("%s log folder is required", LogPrefix)
	}
	if strings.EqualFold(appName, "") {
		return fmt.Errorf("%s app name is required", LogPrefix)
	}

	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s.json", loggerConfig.LogFolder+appName+"-old", "%Y-%m-%d"),
		rotatelogs.WithLinkName(loggerConfig.LogFolder+appName+".json"),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(500),
	)
	if err != nil {
		return fmt.Errorf("%s failed to initialize log file::%s", LogPrefix, err.Error())
	}

	log.SetFormatter(
		&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			FieldMap: log.FieldMap{
				"app": appName,
			},
		})
	// set output level
	if len(loggerConfig.LogLevel) == 0 {
		loggerConfig.LogLevel = "info"
	}
	l, err := log.ParseLevel(loggerConfig.LogLevel)
	if err != nil {
		fmt.Printf(`%s Bad log level %s: %s\n`, LogPrefix, loggerConfig.LogLevel, err.Error())
		return err
	}
	log.SetLevel(l)
	log.SetOutput(writer)
	Log("%s Logger initialized successfully", LogPrefix)
	Log("%sLog folder:%s Log level:%v App Name:%s", LogPrefix, loggerConfig.LogFolder, l, appName)
	return nil
}

func Log(msgFormat string, params ...interface{}) {
	LogInfo(msgFormat, params...)
}

func LogInfo(msgFormat string, params ...interface{}) {
	msg := fmt.Sprintf(msgFormat, params...)
	log.WithField("app", appName).Info(msg)
}

func LogWarn(msgFormat string, params ...interface{}) {
	msg := fmt.Sprintf(msgFormat, params...)
	log.WithField("app", appName).Warn(msg)
}

func LogError(msgFormat string, params ...interface{}) {
	msg := fmt.Sprintf(msgFormat, params...)
	log.WithField("app", appName).Error(msg)
}

func LogDebug(msgFormat string, params ...interface{}) {
	msg := fmt.Sprintf(msgFormat, params...)
	log.WithField("app", appName).Debug(msg)
}
