package tasks

import (
	"github.com/robfig/cron/v3"
	"github.com/seekincs/rate-limit/configs"
	"github.com/seekincs/rate-limit/internal/model"
)

func InitTasks() {
	// Cache a copy first, in case the scheduled task does not start up,
	// resulting in the configuration not being available to read.
	model.CacheRateLimitConfig()
	c := cron.New()
	_, addFuncErr := c.AddFunc(
		configs.RefreshConfigCron,
		model.CacheRateLimitConfig,
	)
	if addFuncErr != nil {
		panic(addFuncErr)
	}
	c.Start()
}
