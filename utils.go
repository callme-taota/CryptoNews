package main

import "time"

func GetCurrentDay() time.Time {
	now := time.Now()
	return now
}

func GetPreDay() time.Time {
	now := time.Now()
	preDay := now.AddDate(0, 0, -1)
	return preDay
}

func TimeCaseCmp(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func JSUnixTimeToTime(jsTimestamp int) time.Time {
	// 将毫秒转换为秒
	sec := int64(jsTimestamp) / 1000
	// 计算余数，即纳秒部分
	nsec := (int64(jsTimestamp) % 1000) * 1e6
	// 使用time.Unix函数将秒和纳秒转换为time.Time类型
	return time.Unix(sec, nsec)
}
