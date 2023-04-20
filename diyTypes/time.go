package diyTypes

import (
	"fmt"
	"github.com/fulldog/utools/timex"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"time"
)

// DiyTime 解析字符串时间
type DiyTime struct {
	time.Time
	Layout string
}

// DonetTime 解析.net的时间 Date\((\d+)-\d+\)
type DonetTime time.Time

// UnixTime 解析时间戳
type UnixTime time.Time

func (t *DiyTime) MarshalJSON() ([]byte, error) {
	var stamp string
	if t.Layout == "" {
		stamp = fmt.Sprintf("\"%s\"", t.Time.String())
	} else {
		stamp = fmt.Sprintf("\"%s\"", t.Time.Format(t.Layout))
	}
	return []byte(stamp), nil
}
func (t *DiyTime) UnmarshalJSON(data []byte) error {
	//先尝试正常转
	err := t.Time.UnmarshalJSON(data)
	if err != nil {
		s := string(data)
		t.Time, err = time.Parse(timex.DateTime, s)
		t.Layout = timex.DateTime
		if err != nil {
			t.Layout = timex.DateOnly
			t.Time, err = time.Parse(timex.DateOnly, s)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

// ToYmdHis 格式化
func (t *DiyTime) ToYmdHis() string {
	return t.Time.Format("2006-01-02 15:04:05")
}

// ToYmd 简单格式化
func (t *DiyTime) ToYmd() string {
	return t.Time.Format("2006-01-02")
}

func (t *DiyTime) ToString() string {
	return t.Time.String()
}

func (dt *DonetTime) UnmarshalJSON(data []byte) error {
	//先尝试正常转
	var t time.Time
	err := t.UnmarshalJSON(data)
	if err == nil {
		*dt = DonetTime(t)
		return nil
	}
	// 从输入字符串中提取时间戳
	s := string(data)
	re := regexp.MustCompile(`Date\((\d+)-\d+\)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return errors.New("时间格式错误" + s)
	}
	timestamp, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return errors.New("时间格式错误" + s)
	}
	// 将时间戳转换为 time.Time 对象
	*dt = DonetTime(time.Unix(timestamp/1000, 0).In(timex.TimeZone))
	return nil
}
func (dt *DonetTime) ToTime() time.Time {
	return time.Time(*dt)
}

func (dt *UnixTime) UnmarshalJSON(data []byte) error {
	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return errors.New("时间格式错误" + string(data))
	}
	// 将时间戳转换为 time.Time 对象
	*dt = UnixTime(time.Unix(timestamp/1000, 0).In(timex.TimeZone))
	return nil
}
func (dt *UnixTime) ToTime() time.Time {
	return time.Time(*dt)
}
func (dt *UnixTime) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(dt.ToTime().Unix())
}
