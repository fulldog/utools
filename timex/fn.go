package timex

import (
	"time"
)

var TimeZone, _ = time.LoadLocation("Asia/Shanghai")
var TimeEnd59 = 86399 * time.Second

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

func GetMonthStartDayStr(dt time.Time, layout string) string {
	return GetMonthStartDay(dt).Format(layout)
}

func GetMonthStartDay(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
}

func GetYMD(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
}

func TimeSpan(day, hour, min, seconds int) int {
	return 86400*day + hour*3600 + min*60 + seconds
}

func Parse(tm string) time.Time {
	if tm == "" {
		return time.Time{}
	}
	t, err := time.ParseInLocation(DateOnly, tm, TimeZone)
	if err == nil {
		return t
	}
	t, err = time.ParseInLocation(DateTime, tm, TimeZone)
	if err == nil {
		return t
	}
	return time.Time{}
}

func ParseAtYmsHis(s string) (t time.Time, err error) {
	t, err = time.ParseInLocation(DateOnly, s, TimeZone)
	if err == nil {
		t = t.Add(86399 * time.Second)
	}
	return
}

// AddMonth 计算月数 修正版
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
