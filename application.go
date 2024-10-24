package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Need 2 input argument, such as [start_utc_time] [end_utc_time]")
		return
	}
	timeSlots, err := splitToBucket(args[1], args[2], 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	size := len(timeSlots)
	hourly = make([]AccumulatorList, size)
	errorChan := make(chan error, size)

	var wg sync.WaitGroup
	for i := 0; i < size; i++ {
		url := toQueryUrl(timeSlots[i].start, timeSlots[i].end)
		hourly[i] = make(AccumulatorList, 0)
		wg.Add(1)
		go asyncSimpleHttpGetCall(url, &hourly[i], process, &wg, errorChan)
	}
	wg.Wait()

	output()
}

func output() {
	for i := 0; i < len(hourly); i++ {
		for j := 0; j < len(hourly[i]); j++ {
			fmt.Printf("%s %8.4f\n", hourly[i][j].hourly, hourly[i][j].GetAverage())
		}
	}
}
