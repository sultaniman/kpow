package cron

import (
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func NewScheduler(schedule string) *cron.Cron {
	return cron.New(cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
	))
}
