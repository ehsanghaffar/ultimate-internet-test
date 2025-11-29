// Package modules provides various internet connectivity tests including HTTP, speed, VPN, and ping checks.
package modules

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ehsanghaffar/ultimate-internet-test/config"
	"github.com/ehsanghaffar/ultimate-internet-test/utils"
)

// TestHTTP performs an HTTP test on the given URL and returns the result.
// It accepts a config parameter for timeout configuration and returns an HTTPTest result with any error encountered.
// The function logs all HTTP response details including status, TLS information, and headers.
//
// Parameters:
//   - url: The URL to test (HTTP or HTTPS)
//   - cfg: Configuration containing timeout settings
//
// Returns:
//   - *HTTPTest: Pointer to HTTPTest struct containing the test results and any errors
//
// Example:
//
//	cfg := config.New()
//	result := TestHTTP("https://example.com", cfg)
//	if result.Error != "" {
//	    log.Println("Test failed:", result.Error)
//	}
func TestHTTP(url string, cfg *config.Config) *utils.HTTPTest {
	result := &utils.HTTPTest{
		URL: url,
	}

	log.Println("URL:", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		result.Error = err.Error()
		log.Println("Error creating request:", url, err)
		return result
	}

	// Use provided timeout from config
	client := http.Client{
		Timeout: cfg.HTTPTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Println("Redirect:", req.URL)
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Error = err.Error()
		log.Println("Error sending request:", url, err)
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.Status
	result.Proto = resp.Proto

	log.Println("Response status:", resp.Status, resp.Proto)

	if resp.TLS != nil {
		result.TLSVersion = fmt.Sprintf("%d", resp.TLS.Version)
		result.CipherSuite = fmt.Sprintf("%d", resp.TLS.CipherSuite)
		result.ServerName = resp.TLS.ServerName

		log.Println("Response TLS version:", resp.TLS.Version)
		log.Println("Response TLS cipher suite:", resp.TLS.CipherSuite)
		log.Println("Response TLS server name:", resp.TLS.ServerName)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		log.Println("Error reading response:", url, err)
		return result
	}

	result.ResponseLength = len(body)

	for k, v := range resp.Header {
		log.Println("Response header:", k, v)
	}

	log.Println("Response length:", len(body))
	fmt.Println("------------------------------------------------------------")

	return result
}
