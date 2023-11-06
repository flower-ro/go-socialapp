package notlogin

import (
	"github.com/marmotedu/iam/pkg/log"
	utils "go-socialapp/internal/pkg/util"
	"time"
)

var defaultInterval = 90 * time.Second

// 创建时间超过1分钟 就删除
func (t tmpWaClientCache) scan() {
	interval := defaultInterval
	for {
		pre := len(t.tmpWaClients)
		interval = t.delExpire(interval)
		log.Infof("delete tmp client ,now is %d,before is %d,interval %v", len(t.tmpWaClients), pre, interval)
		time.Sleep(interval)

	}

}

func (t *tmpWaClientCache) delExpire(expireTime time.Duration) time.Duration {
	var minInterval time.Duration
	var now = utils.GetCurrentTime()
	for fileName, client := range t.tmpWaClients {
		tmpClient := client
		lastHead := tmpClient.CreateTime
		expectExpireTime := lastHead.Add(expireTime)
		if now.After(expectExpireTime) || now.Equal(expectExpireTime) {
			tmpClient.WaCli.Disconnect()
			delete(t.tmpWaClients, fileName)
		} else {
			expectExpireTime = lastHead.Add(defaultInterval)
			interval := expectExpireTime.Sub(now)
			if minInterval == 0 {
				minInterval = interval
			}
			if interval < minInterval {
				minInterval = interval
			}
		}
	}
	if minInterval == 0 {
		minInterval = defaultInterval
	}
	return minInterval
}
