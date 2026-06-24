package scheduler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/k3nsei/schedulord/internal/config"
	gh "github.com/k3nsei/schedulord/internal/github"
	"github.com/k3nsei/schedulord/internal/handlers"
	"github.com/k3nsei/schedulord/internal/models"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron        *cron.Cron
	location    *time.Location
	githubToken string
}

func New(loc *time.Location, token string) *Scheduler {
	return &Scheduler{
		cron:        cron.New(cron.WithLocation(loc)),
		location:    loc,
		githubToken: token,
	}
}

func (s *Scheduler) Start(cfg *config.Config) error {
	s.cron.Start()

	for i, ts := range cfg.Tasks {
		t := ts.Task

		handler, ok := handlers.ResolveHandler(t.GetBase().Type)
		if !ok {
			slog.Error(
				"Task handler not found",
				"taskID", i,
				"taskType", t.GetBase().Type,
			)
			continue
		}

		s.register(i, t, handler)
	}

	return nil
}

func (s *Scheduler) register(
	idx int,
	task models.Task,
	handler models.TaskHandler,
) {
	base := task.GetBase()

	if base.Disabled {
		slog.Info(
			"Task skipped (disabled)",
			"taskID", idx,
			"taskType", base.Type,
		)
		return
	}

	entryID, err := s.cron.AddFunc(base.Cron, func() {
		s.execute(idx, task, handler)
	})
	if err != nil {
		slog.Error(
			"Task registration failed",
			"taskID", idx,
			"taskType", base.Type,
			"error", err,
		)
		return
	}

	entry := s.cron.Entry(entryID)

	slog.Info(
		"Task scheduled",
		"taskID", idx,
		"taskType", base.Type,
		"nextRun", entry.Next,
	)
}

func (s *Scheduler) execute(
	idx int,
	task models.Task,
	handler models.TaskHandler,
) {
	base := task.GetBase()
	startedAt := time.Now()

	slog.Info(
		"Task started",
		"taskID", idx,
		"taskType", base.Type,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx = gh.WithGithubToken(ctx, s.githubToken)

	err := handler.Run(ctx, base.Payload)
	duration := time.Since(startedAt)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		slog.Error(
			"Task timed out",
			"taskID", idx,
			"taskType", base.Type,
			"duration", duration,
		)
	case err != nil:
		slog.Error(
			"Task failed",
			"taskID", idx,
			"taskType", base.Type,
			"duration", duration,
			"error", err,
		)
	default:
		slog.Info(
			"Task completed",
			"taskID", idx,
			"taskType", base.Type,
			"duration", duration,
		)
	}
}
