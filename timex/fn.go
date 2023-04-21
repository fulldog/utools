package timex

import (
	"time"
)

// CompareTime 比较时间大小
func CompareTime(t1, t2 time.Time, cond string) bool {
	switch cond {
	case ">":
		return t1.After(t2)
	case "<":
		return t1.Before(t2)
	case "=":
		return t1.Equal(t2)
	case ">=":
		return t1.Equal(t2) || t1.After(t2)
	case "<=":
		return t1.Equal(t2) || t1.Before(t2)
	}
	return false
}

// GetMonthStartDayStr 返回每月的第一天
func GetMonthStartDayStr(dt time.Time, layout string) string {
	return GetMonthStartDay(dt).Format(layout)
}

// GetMonthStartDay 返回每月的第一天
func GetMonthStartDay(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
}

// GetYmdStart 当天的开始时间 00:00:00
func GetYmdStart(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
}

// GetYmdEnd 当天的开始时间 23:59:59
func GetYmdEnd(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 23, 59, 59, 0, dt.Location())
}

// TimeSpan 计算秒数
func TimeSpan(day, hour, min, seconds int) int {
	return 86400*day + hour*3600 + min*60 + seconds
}

// ParseYmd 格式化时间
func ParseYmd(tm string) (time.Time, error) {
	return time.ParseInLocation(DateOnly, tm, time.Local)
}

// ParseYmdHis 格式化时间
func ParseYmdHis(tm string) (time.Time, error) {
	return time.ParseInLocation(DateTime, tm, time.Local)
}

// AddMonth 计算月数 当后退时，自动修正到该月份的最后一天
func AddMonth(t time.Time, sub int) time.Time {
	if sub == 0 {
		return t
	}
	if sub > 0 {
		return t.AddDate(0, sub, 0)
	}

	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	maxDay := time.Date(year, month, 0, 0, 0, 0, 0, t.Location()).Day()
	if day > maxDay {
		day = maxDay
	}
	return time.Date(year, month+time.Month(sub), day, hour, min, sec, 0, t.Location())
}

// Time2UtcForMsql 默认会全局使用cts时间类型，但是SqlServer使用utc 会导致查询数据异常，需要转成utc时间
func Time2UtcForMsql(tm time.Time) time.Time {
	tm, _ = time.Parse(DateTime, tm.Format(DateTime))
	return tm
}

func GetTicker(second time.Duration) *time.Ticker {
	return time.NewTicker(second * time.Second)
}

func NewTicket(fn func(), t *time.Ticker) {
	for {
		fn()
		<-t.C
	}
}
