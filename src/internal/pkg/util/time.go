package utils

import "time"

func GetCurrentTime() time.Time {
	now := time.Now() //使用的 CST(中国标准时间)
	now.Format("2006-01-02 15:04:05")
	return now
}
