package main

import (
	"math"
	"sync"
)

// Accumulator struct definition
type Accumulator struct {
	hourly            string
	sampleCount       int
	accumulatedValues float64
	mu                sync.Mutex // Mutex for concurrent access
}

// NewAccumulator is a constructor for creating an Accumulator
func NewAccumulator(hourly string) *Accumulator {
	return &Accumulator{
		hourly:            hourly,
		sampleCount:       0,
		accumulatedValues: 0.0,
	}
}

// GetAverage calculates the average of accumulated values
func (a *Accumulator) GetAverage() float64 {
	if a.sampleCount == 0 {
		return 0.0 // Avoid division by zero
	}
	value := a.accumulatedValues / float64(a.sampleCount)
	return math.Round(value*10000) / 10000.0
}

// AddSample adds a sample value to the Accumulator
func (a *Accumulator) AddSample(value float64) {
	a.mu.Lock()         // Lock for concurrent access
	defer a.mu.Unlock() // Unlock at the end

	a.accumulatedValues += value
	a.sampleCount++
}

type AccumulatorList []Accumulator
