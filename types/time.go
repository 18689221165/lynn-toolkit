package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time is alias type for time.Time
type Time time.Time

const (
	TimeFormart = "2006-01-02 15:04:05"
	zone        = "Asia/Shanghai"
)

func NowTime() Time {
	return Time(time.Now())
}

func (t *Time) Format(layout string) string {
	return time.Time(*t).Format(layout)
}

func ParseTimeStr(layout, value string) (Time, error) {
	t, err := time.ParseInLocation(layout, value, time.Local)
	return Time(t), err
}

func (t Time) IsZero() bool {
	return t.IsZero()
}

// UnmarshalJSON implements json unmarshal interface.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

// MarshalJSON implements json marshal interface.
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeFormart)
}

func (t Time) local() time.Time {
	loc, _ := time.LoadLocation(zone)
	return time.Time(t).In(loc)
}

// Value ...
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time 注意是指针类型 method
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t Time) AtTime(timeStr string) time.Time {
	att, _ := time.ParseInLocation(TimeFormart, t.Format("2006-01-02")+" "+timeStr, time.Local)
	return att
}
