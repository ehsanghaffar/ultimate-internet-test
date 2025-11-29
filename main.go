package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/ehsanghaffar/ultimate-internet-test/config"
	"github.com/ehsanghaffar/ultimate-internet-test/modules"
	"github.com/ehsanghaffar/ultimate-internet-test/utils"
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Initialize configuration with defaults
	cfg := config.New()

	// Parse command-line arguments for custom URLs
	args := flag.Args()
	if len(args) > 0 {
		runHTTPTests(args, cfg)
		return
	}

	// Run all default tests
	runAllTests(cfg)
}

// runHTTPTests runs HTTP tests on the provided URLs
func runHTTPTests(urls []string, cfg *config.Config) {
	var wg sync.WaitGroup
	httpTests := make([]*utils.HTTPTest, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()
			httpTests[index] = modules.TestHTTP(u, cfg)
		}(i, url)
	}

	wg.Wait()

	// Convert to non-pointer slice for storage
	var results []utils.HTTPTest
	for _, test := range httpTests {
		if test != nil {
			results = append(results, *test)
		}
	}

	// Save results to file
	testResults := &utils.TestResults{
		HTTPTests: results,
	}

	if err := utils.SaveResults(testResults, cfg.ResultsFilePath, config.FilePermissions); err != nil {
		log.Printf("Error saving results: %v\n", err)
	}
}

// runAllTests runs all available tests concurrently
func runAllTests(cfg *config.Config) {
	var wg sync.WaitGroup

	// Initialize result containers
	var (
		httpTests  []*utils.HTTPTest
		speedTests []*utils.SpeedTest
		vpnTest    *utils.VPNTest
		pingTest   *utils.PingTest
		mu         sync.Mutex
	)

	// Default HTTP test URLs
	httpURLs := []string{
		"http://www.google.com/",
		"https://www.google.com/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://leader.ir/",
	}

	// Default speed test URLs
	speedURLs := []string{
		"https://ehsanghaffarii.ir",
		"https://google.com",
	}

	// Run HTTP tests concurrently
	for _, url := range httpURLs {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := modules.TestHTTP(u, cfg)
			mu.Lock()
			httpTests = append(httpTests, result)
			mu.Unlock()
		}(url)
	}

	// Run speed tests concurrently
	for _, url := range speedURLs {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := modules.CheckSpeed(u, cfg)
			mu.Lock()
			speedTests = append(speedTests, result)
			mu.Unlock()
		}(url)
	}

	// Run VPN check (sequential, as it involves IP detection)
	wg.Add(1)
	go func() {
		defer wg.Done()
		vpnTest = modules.CheckVPN("http://checkip.dyndns.org/")
	}()

	// Run ping tests concurrently
	pingDomains := []string{
		"www.ehsanghaffarii.ir",
		"www.google.com",
	}

	for _, domain := range pingDomains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			result := modules.PingCheck(d, cfg)
			mu.Lock()
			if pingTest == nil {
				pingTest = result
			}
			mu.Unlock()
		}(domain)
	}

	// Wait for all tests to complete
	wg.Wait()

	// Convert pointers to values for storage
	var httpTestsValues []utils.HTTPTest
	for _, test := range httpTests {
		if test != nil {
			httpTestsValues = append(httpTestsValues, *test)
		}
	}

	var speedTestsValues []utils.SpeedTest
	for _, test := range speedTests {
		if test != nil {
			speedTestsValues = append(speedTestsValues, *test)
		}
	}

	// Create aggregated results
	testResults := &utils.TestResults{
		HTTPTests:  httpTestsValues,
		SpeedTests: speedTestsValues,
	}

	if vpnTest != nil {
		testResults.VPNTest = *vpnTest
	}

	if pingTest != nil {
		testResults.PingTest = *pingTest
	}

	// Save all results at once
	if err := utils.SaveResults(testResults, cfg.ResultsFilePath, config.FilePermissions); err != nil {
		log.Printf("Error saving results: %v\n", err)
	} else {
		fmt.Printf("Results saved to %s\n", cfg.ResultsFilePath)
	}
}
