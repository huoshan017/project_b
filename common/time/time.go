package time

import (
	"time"
)

type Duration time.Duration

const (
	Nanosecond  = Duration(time.Nanosecond)
	Microsecond = Duration(time.Microsecond)
	Millisecond = Duration(time.Millisecond)
	Second      = Duration(time.Second)
	Minute      = Duration(time.Minute)
	Hour        = Duration(time.Hour)
)

type CustomTime struct {
	time.Time
}

func (t CustomTime) Sub(tt CustomTime) Duration {
	return Duration(t.Time.Sub(tt.Time))
}

func (t CustomTime) Add(duration Duration) CustomTime {
	return CustomTime{t.Time.Add(time.Duration(duration))}
}

func (t CustomTime) AddDate(years, months, days int) CustomTime {
	return CustomTime{t.Time.AddDate(years, months, days)}
}

func Now() CustomTime {
	return CustomTime{time.Now()}
}

func Since(t CustomTime) Duration {
	return Duration(time.Since(t.Time))
}

func Until(t CustomTime) Duration {
	return Duration(time.Until(t.Time))
}
