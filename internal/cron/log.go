package cron

import (
	"newblog/internal/config"
	"newblog/internal/logger"
	"os"
	"time"
)

type Log struct{}

func (l *Log) GetRetryTimes() int {
	return 3
}

func (l *Log) Exec() error {
	// 切割日志每天使用新文件
	logger.Init()

	// 扫描日志目录
	files, err := os.ReadDir(config.Global.Log.Path)
	if err != nil {
		return err
	}

	// 获取15天前的时间
	thirtyDaysAgo := time.Now().AddDate(0, 0, -15)

	// 遍历日志文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 获取文件信息
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		// 如果文件修改时间在30天前，则删除
		if fileInfo.ModTime().Before(thirtyDaysAgo) {
			os.Remove(config.Global.Log.Path + file.Name())
		}
	}

	return nil
}
