package drop

import (
	"context"
	"log/slog"

	"github.com/go-co-op/gocron/v2"
)

func RegisterJobs(scheduler gocron.Scheduler, service DropService) (gocron.Job, error) {
	return scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(0, 0, 0),
			),
		),
		gocron.NewTask(
			func() {
				err := service.DeleteExpiredDrops(context.TODO())
				if err != nil {
					slog.Error("failed deleting expired drops", "error", err.Error())
				} else {
					slog.Info("successfully deleted expired drops")
				}
			},
		),
	)
}
