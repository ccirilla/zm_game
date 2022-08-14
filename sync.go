package main

import (
	"time"
)

func DocSyncTaskCronJob() {
	last_int := 0
	now_int := 0
	ticker := time.NewTicker(time.Minute * 1) // 每分钟执行一次
	for range ticker.C {
		now_int = int(time.Now().Unix())
		if time.Now().Hour() == 2 && (now_int - last_int) > 24*60*60{
			last_int = now_int - 30 * 60
			SyncAllDevice()
		}
	}
}