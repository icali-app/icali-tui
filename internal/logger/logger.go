package logger

import (
	"context"
	appConf "github.com/icali-app/icali-tui/internal/config"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		conf := appConf.Get()

		loglevel, err := zerolog.ParseLevel(conf.Logging.LogLevel)
		if err != nil {
			loglevel = zerolog.DebugLevel // default to debug
		}

		err = os.MkdirAll(conf.Logging.LogDir, 0744)
		if err != nil {
			return
		}

		logFilePath := filepath.Join(
			conf.Logging.LogDir,
			"icali-tui.log.json")

		logfile, err := os.OpenFile(
			logFilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)

		if err != nil {
			panic(err)
		}

		log = zerolog.New(logfile).
			Level(loglevel).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger().
			With().
			Caller().
			Logger()
	})

	return log
}

func GetCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
