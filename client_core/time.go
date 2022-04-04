package client_core

import (
	custom_time "project_b/common/time"
)

const (
	defaultTotalTimeSyncCount = 40 // 缺省时间同步总次数
)

var (
	ctime              *custom_time.Client  = custom_time.NewClient() // 时间客户端
	c2sDeltaTime       custom_time.Duration                           // 客户端与服务器的时间差
	totalTimeSyncCount int                                            // 时间同步总计数
	timeSyncCount      int                                            // 时间同步计数
)

func SetTimeSyncTotalCount(c int) {
	totalTimeSyncCount = c
}

func SetSyncSendTime(t custom_time.CustomTime) {
	ctime.SetSendTime(t)
}

func SetSyncRecvAndServerTime(recvTime, serverTime custom_time.CustomTime) {
	ctime.SetRecvAndServerTime(recvTime, serverTime)
	timeSyncCount += 1
}

func GetSyncSendTime() custom_time.CustomTime {
	return ctime.GetSendTime()
}

func IsTimeSyncEnd() bool {
	if totalTimeSyncCount == 0 {
		totalTimeSyncCount = defaultTotalTimeSyncCount
	}
	return timeSyncCount >= totalTimeSyncCount
}

func GetNetworkDelay() custom_time.Duration {
	return ctime.GetDelay()
}

func GetNetworkAvgDelay() custom_time.Duration {
	return ctime.GetAvgDelay()
}

func GetSyncServTime() custom_time.CustomTime {
	delay := ctime.GetDelay()
	if delay < 0 {
		panic("network delay invalid")
	}
	if c2sDeltaTime == 0 {
		var b bool
		c2sDeltaTime, b = ctime.GetDeltaC2S()
		if !b {
			panic("GetDeltaC2S failed")
		}
	}
	return custom_time.Now().Add(-c2sDeltaTime + ctime.GetDelay())
}
