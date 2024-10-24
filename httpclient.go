package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func asyncSimpleHttpGetCall(
	url string,
	buckets *AccumulatorList,
	callback Callback,
	wg *sync.WaitGroup,
	errorChan chan<- error,
) {
	defer wg.Done()

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errorChan <- err
		return
	}

	// Optionally, set headers (if needed)
	req.Header.Set("Accept", "application/json")

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		errorChan <- err
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			errorChan <- err
			return
		}
	}(resp.Body) // Ensure the response body is closed after reading

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		errorChan <- fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
		return
	}

	// Create a new buffered reader
	scanner := bufio.NewScanner(resp.Body)

	// Read the response body line by line
	for scanner.Scan() {
		// Get the current line
		line := scanner.Text()
		// Process the line (print it, for example)
		err := callback(line, buckets)
		if err != nil {
			errorChan <- err
			return
		}
	}

	// Check for errors that may have occurred during scanning
	if err := scanner.Err(); err != nil {
		errorChan <- err
		return
	}
}

type Callback func(string, *AccumulatorList) error
