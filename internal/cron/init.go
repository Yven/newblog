package cron

import (
	"fmt"
	"log"
	"newblog/internal/global"
	"time"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	c         *cron.Cron
	entries   []cron.Entry
	missonMap map[ExecFunc]string
	logger    *SlogAdapter
}

func NewCronService() *CronService {
	clone := *global.Logger
	logger := &SlogAdapter{logger: &clone}
	return &CronService{
		c: cron.New(
			cron.WithLogger(logger),
			cron.WithChain(cron.Recover(logger)),
		),
		logger: logger,
	}
}

type ExecFunc interface {
	Exec() error
	GetRetryTimes() int
}

func (s *CronService) Register() {
	// Field name   | Mandatory? | Allowed values  | Allowed special characters
	// ----------   | ---------- | --------------  | --------------------------
	// Minutes      | Yes        | 0-59            | * / , -
	// Hours        | Yes        | 0-23            | * / , -
	// Day of month | Yes        | 1-31            | * / , - ?
	// Month        | Yes        | 1-12 or JAN-DEC | * / , -
	// Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
	s.missonMap = map[ExecFunc]string{
		// 每天 4 点钟同步一次 notion 文章
		&NotionBlog{}: "0 4 * * *",
		// 每天 00:00 分割日志
		&Log{}: "0 0 * * *",
	}

	for mission, spec := range s.missonMap {
		id, _ := s.c.AddFunc(spec, func() {
			// 重试
			for i := 0; i < mission.GetRetryTimes(); i++ {
				if i > 0 {
					time.Sleep(5 * time.Second)
				}
				log.Printf("定时任务开始执行[第%d次]: %T", i+1, mission)

				err := mission.Exec()
				if err == nil {
					break
				}

				s.logger.Error(err, "定时任务执行失败", "mission", fmt.Sprintf("%T", mission), "retry_times", i+1)
			}
		})
		s.entries = append(s.entries, s.c.Entry(id))
	}

	s.c.Start()
}

func (s *CronService) Stop() {
	ctx := s.c.Stop()
	<-ctx.Done()
	fmt.Println("定时任务已停止")
}
