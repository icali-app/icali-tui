package logger

import (
	"context"
	"os"
	"runtime/debug"
	"strconv"
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

		logLevelInt, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		
		logLevel := int8(logLevelInt)
		if int(logLevel) != logLevelInt {
			logLevel = int8(zerolog.InfoLevel)
		}

		if err != nil {
			logLevel = int8(zerolog.DebugLevel) // default to DEBUG
		}

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

		logfile, err := os.OpenFile(
			"icali-tui.log",
			os.O_APPEND | os.O_CREATE | os.O_WRONLY,
			0664,
		)

		if err != nil {
			panic(err)
		}

		log = zerolog.New(logfile).
			Level(zerolog.Level(logLevel)).
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

