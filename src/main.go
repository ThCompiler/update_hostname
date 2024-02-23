package main

import (
	"github.com/ThCompiler/go.beget.api/core"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"update_hostname/config"
	"update_hostname/internal/logger"
	"update_hostname/internal/updator"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	l, deferFunc := PrepareLogger(cfg.LoggerInfo)
	defer func() {
		deferFunc()
		_ = l.Sync()
	}()

	updater := updator.NewUpdater(
		core.Client{
			Login:    cfg.Login,
			Password: cfg.Password,
		},
		l,
		cfg.Domain,
	)

	loop(updater, l, cfg.UpdateHours)
}

func loop(updater *updator.Updater, l logger.Interface, updateHour int64) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(updateHour) * time.Hour)
	updater.Update()
loop:
	for {
		select {
		case s := <-interrupt:
			l.Info(s.String())
			break loop
		case <-ticker.C:
			updater.Update()
		}
	}
}

func ConvertConfigLoggerParam(cfg config.LoggerInfo) logger.Params {
	return logger.Params{
		AppName:                  cfg.AppName,
		LogDir:                   cfg.Directory,
		Level:                    cfg.Level,
		UseStdAndFile:            cfg.UseStdAndFile,
		AddLowPriorityLevelToCmd: cfg.AllowShowLowLevel,
	}
}

func PrepareLogger(cfg config.LoggerInfo) (l *logger.Logger, deferFunc func()) {
	var logOut io.Writer

	if cfg.Directory != "" {
		file, err := logger.OpenLogDir(cfg.Directory)
		if err != nil {
			log.Fatalf("create logger error: %s", err)
		}

		deferFunc = func() {
			err = file.Close()
			log.Fatalf("close log file error: %s", err)
		}

		logOut = file
	} else {
		logOut = os.Stderr
	}

	l = logger.New(ConvertConfigLoggerParam(cfg), logOut)

	return l, deferFunc
}
