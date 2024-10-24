package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	hourly []AccumulatorList
	mu     sync.Mutex
)

func toQueryUrl(start string, end string) string {
	return fmt.Sprintf("https://tsserv.tinkermode.dev/data?begin=%s&end=%s", start, end)
}

func process(token string, buckets *AccumulatorList) error {
	err := aggregateToken(token, buckets)
	if err != nil {
		return err
	}
	return nil
}

func aggregateToken(input string, buckets *AccumulatorList) error {
	// Lock for concurrent safety
	mu.Lock()
	defer mu.Unlock()

	// Split the input string by spaces
	token := strings.Fields(input)
	if len(token) < 2 {
		return fmt.Errorf("input does not contain expected tokens")
	}

	// Check the timeslot using the first token
	timeslot, err := toHourlyUTCTime(token[0])
	if err != nil {
		return err
	}
	currSize := len((*buckets))
	latest := ""
	if currSize > 0 {
		latest = (*buckets)[len((*buckets))-1].hourly
	}
	if timeslot != "" {
		// Check if the timeslot is different from the latest
		if latest != timeslot {
			(*buckets) = append((*buckets), *NewAccumulator(timeslot))
			latest = timeslot
		}

		value, err := strconv.ParseFloat(token[1], 64)
		if err != nil {
			return err
		}
		(*buckets)[len((*buckets))-1].AddSample(value)
	}
	return nil
}
