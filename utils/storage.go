package utils

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

var (
	// resultsMutex protects concurrent access to the results file
	resultsMutex sync.Mutex
)

// LoadResults loads test results from a JSON file
func LoadResults(filePath string) (*TestResults, error) {
	resultsMutex.Lock()
	defer resultsMutex.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return empty results
			return &TestResults{
				Timestamp: time.Now(),
			}, nil
		}
		return nil, NewNetworkError("Storage", "failed to read results file", err)
	}

	var results TestResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, NewParseError("Storage", "failed to parse results JSON", err)
	}

	return &results, nil
}

// SaveResults saves test results to a JSON file
func SaveResults(results *TestResults, filePath string, filePermissions os.FileMode) error {
	if results == nil {
		return NewValidationError("Storage", "results cannot be nil")
	}

	resultsMutex.Lock()
	defer resultsMutex.Unlock()

	// Set timestamp if not already set
	if results.Timestamp.IsZero() {
		results.Timestamp = time.Now()
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return NewParseError("Storage", "failed to marshal results to JSON", err)
	}

	if err := os.WriteFile(filePath, data, filePermissions); err != nil {
		return NewNetworkError("Storage", "failed to write results file", err)
	}

	return nil
}

// AppendResult appends a single result to the existing results and saves
func AppendResult(httpTests []HTTPTest, speedTests []SpeedTest, vpnTest *VPNTest, pingTest *PingTest, filePath string, filePermissions os.FileMode) error {
	// Load existing results
	results, err := LoadResults(filePath)
	if err != nil {
		// If file doesn't exist, create new results
		results = &TestResults{
			Timestamp: time.Now(),
		}
	}

	// Append new results
	if len(httpTests) > 0 {
		results.HTTPTests = append(results.HTTPTests, httpTests...)
	}
	if len(speedTests) > 0 {
		results.SpeedTests = append(results.SpeedTests, speedTests...)
	}
	if vpnTest != nil {
		results.VPNTest = *vpnTest
	}
	if pingTest != nil {
		results.PingTest = *pingTest
	}

	results.Timestamp = time.Now()

	return SaveResults(results, filePath, filePermissions)
}
