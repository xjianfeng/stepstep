package utils

import (
	"math"
	"time"
)

var (
	sinTimeSec, _ = time.ParseInLocation("2006-01-02 15:04:05", "2010-01-01 00:00:00", time.Local)
)

func GetRealDayNo() int {
	DayNo := int(math.Ceil(time.Now().Sub(sinTimeSec).Hours() / 24))
	return DayNo
}
