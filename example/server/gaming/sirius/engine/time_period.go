package engine

import (
	"time"
)

// 表示一段周期性增长的时间段
// [Start + n * Interval, Start + n * Interval + Duration)
// 当Duration == 0 时退化为一个周期性的时间点
type TimePeriod struct {
	Start    time.Time
	Interval time.Duration
	Duration time.Duration
}

func (tp TimePeriod) Empty() bool {
	return tp.Start.IsZero()
}

func (tp TimePeriod) NextTimeFromNow() time.Time {
	return tp.NextTime(time.Now())
}

func (tp TimePeriod) NextTime(from time.Time) time.Time {
	if tp.Start.After(from) {
		return tp.Start
	}
	if tp.Interval == 0 {
		return from
	}

	interval := from.Sub(tp.Start)
	n := interval / tp.Interval
	next := tp.Start.Add(n * tp.Interval)
	for !next.After(from) {
		next = next.Add(tp.Interval)
	}

	return next
}

func (tp TimePeriod) Contains(t time.Time) bool {
	if tp.Interval == 0 {
		return false
	}
	if tp.Start.After(t) {
		return false
	}
	// Pre-condition: tp.Start <= t
	next := tp.Start.Add(tp.Interval)

	// Post-condition: start <= t && t < next
	start := tp.Start
	for next.Before(t) || next.Equal(t) {
		start.Add(tp.Interval)
		next = next.Add(tp.Interval)
	}

	return t.Before(start.Add(tp.Duration))
}
