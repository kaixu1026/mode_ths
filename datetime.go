package main

import (
	"time"
)

type TimeSlot struct {
	start string
	end   string
}

// convert to the standard utc time input from long millisecond timestamp
func toUTCTimeString(timestamp int64) string {
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1_000_000
	t := time.Unix(seconds, nanoseconds).UTC()
	return t.Format(time.RFC3339)
}

// convert to long millisecond timestamp from the standard utc time input
func toMsTimestamp(utcTimeString string) (int64, error) {
	instant, err := time.Parse(time.RFC3339, utcTimeString)
	if err != nil {
		return 0, err
	} else {
		return instant.UnixNano() / int64(time.Millisecond), nil
	}
}

// convert to hourly UTC Time from the standard utc time input
func toHourlyUTCTime(utcTimeString string) (string, error) {
	_, err := toMsTimestamp(utcTimeString)
	if err != nil {
		return "", err
	}
	return utcTimeString[:13] + ":00:00Z", nil
}

func previousSecond(utcTimeString string) (string, error) {
	t, err := time.Parse(time.RFC3339, utcTimeString)
	if err != nil {
		return "", err
	}
	oneSecondEarlier := t.Add(-time.Second).UTC()
	return oneSecondEarlier.Format(time.RFC3339), nil
}

func nextHour(utcTimeString string, num int) (string, error) {
	t, err := time.Parse(time.RFC3339, utcTimeString)
	if err != nil {
		return "", err
	}
	oneSecondEarlier := t.Add(time.Duration(num) * time.Hour).UTC()
	return oneSecondEarlier.Format(time.RFC3339), nil
}

func hourDiff(time1 int64, time2 int64) int {
	return int((time2 - time1) / (1000 * 60 * 60))
}

func nexHour(time int64, num int) int64 {
	return time + (1000*60*60)*(int64(num))
}

func splitToBucket(start string, end string, hours int) ([]TimeSlot, error) {
	ss, err := toHourlyUTCTime(start)
	if err != nil {
		return nil, err
	}
	es, err := toHourlyUTCTime(end)
	if err != nil {
		return nil, err
	}
	nh, err := nextHour(es, 1)
	if err != nil {
		return nil, err
	}

	s, err := toMsTimestamp(ss)
	if err != nil {
		return nil, err
	}
	e, err := toMsTimestamp(nh)
	if err != nil {
		return nil, err
	}

	diff := hourDiff(s, e)
	size := (diff + (hours - 1)) / hours
	slots := make([]TimeSlot, size)
	for i := 0; i < size; i++ {
		p := nexHour(s, hours)
		ss := toUTCTimeString(s)
		ps := toUTCTimeString(p)
		es := toUTCTimeString(e)
		if p > e {
			previousSecond, err := previousSecond(es)
			if err != nil {
				return nil, err
			}
			slots[i] = TimeSlot{ss, previousSecond}
		} else {
			previousSecond, err := previousSecond(ps)
			if err != nil {
				return nil, err
			}
			slots[i] = TimeSlot{ss, previousSecond}
		}
		s = p
	}
	return slots, nil
}
