package cron

import (
	"commmunity/app/zlog"
	"context"
	"sync"
	"time"
)

type CronManager struct {
	ticker *time.Ticker
	stop   chan struct{}
	wg     sync.WaitGroup
}

func NewCronManager(duration time.Duration) *CronManager {
	return &CronManager{
		ticker: time.NewTicker(duration),
		stop:   make(chan struct{}),
	}
}

func (c *CronManager) Start(ctx context.Context, task func(ctx context.Context)) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		zlog.Info("定时任务管理器启动")
		for {
			select {
			case <-c.ticker.C:
				task(ctx)
			case <-c.stop:
				zlog.Info("定时任务管理器停止")
				return
			case <-ctx.Done():
				zlog.Info("定时任务停止")
				return
			}
		}
	}()
}

func (c *CronManager) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	close(c.stop)
	c.wg.Wait()
}
