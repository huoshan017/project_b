package utils

import "fmt"

func FormatSeconds2TimeString(ms uint32) string {
	if ms < 1000 {
		return fmt.Sprintf("00:00.%.2d", ms/10)
	} else if ms >= 1000 && ms < 60*1000 {
		s := ms / 1000
		hms := (ms % 1000) / 10
		return fmt.Sprintf("00:%.2d.%.2d", s, hms)
	} else if ms >= 60*1000 && ms < 3600*1000 {
		m := ms / (60 * 1000)
		s := (ms % (60 * 1000)) / 1000
		hms := (ms % (60 * 1000) % 1000) / 10
		return fmt.Sprintf("%.2d:%.2d.%.2d", m, s, hms)
	}
	return "--:--.--"
}
