package modules

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ehsanghaffar/ultimate-internet-test/config"
	"github.com/ehsanghaffar/ultimate-internet-test/utils"
)

// CheckSpeed performs a speed test by downloading from the given URL and returns the speed result.
// It uses timeout configuration from the config parameter. The function measures download speed in Mbps.
//
// Parameters:
//   - url: The URL to download from for speed testing
//   - cfg: Configuration containing timeout settings
//
// Returns:
//   - *SpeedTest: Pointer to SpeedTest struct containing speed metrics and any errors
//
// Example:
//
//	cfg := config.New()
//	result := CheckSpeed("https://example.com/largefile", cfg)
//	if result.Error == "" {
//	    log.Printf("Download speed: %.2f Mbps\n", result.DownloadMbps)
//	}
func CheckSpeed(url string, cfg *config.Config) *utils.SpeedTest {
	result := &utils.SpeedTest{
		URL: url,
	}

	startTime := time.Now()

	// Create a client with timeout from config
	client := &http.Client{
		Timeout: cfg.SpeedTestTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		result.Error = err.Error()
		log.Println(err)
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		log.Println(err)
		return result
	}

	elapsedTime := time.Since(startTime)
	result.ElapsedTime = elapsedTime
	result.BytesReceived = len(body)

	// Calculate speed in Mbps
	speed := float64(len(body)) / elapsedTime.Seconds()
	result.DownloadMbps = (speed / float64(config.BytesToMegabytes)) * float64(config.BytesToBits)

	log.Println("URL:", url)
	log.Printf("Download speed: %.2f Mbps\n", result.DownloadMbps)
	log.Printf("Elapsed time: %s\n", elapsedTime)
	fmt.Println("------------------------------------------------------------")

	return result
}
