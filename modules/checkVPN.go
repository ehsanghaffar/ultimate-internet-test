package modules

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/ehsanghaffar/ultimate-internet-test/utils"
)

const (
	// IPPattern is the regex pattern to extract IP addresses from HTTP responses
	IPPattern = `IP Address: (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
)

// CheckVPN detects if a VPN or proxy is being used by comparing external and local IPs.
// It makes an HTTP request to an IP detection service and compares the returned IP with localhost.
//
// Parameters:
//   - ipChecker: URL of an IP detection service (e.g., "http://checkip.dyndns.org/")
//
// Returns:
//   - *VPNTest: Pointer to VPNTest struct containing detection status and any errors
//
// Example:
//
//	result := CheckVPN("http://checkip.dyndns.org/")
//	if result.Error == "" {
//	    log.Println("VPN Status:", result.Status)
//	}
func CheckVPN(ipChecker string) *utils.VPNTest {
	result := &utils.VPNTest{}

	resp, err := http.Get(ipChecker)
	if err != nil {
		result.Error = err.Error()
		fmt.Println(err)
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		fmt.Println(err)
		return result
	}

	// Parse external IP with bounds checking
	re := regexp.MustCompile(IPPattern)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		result.Error = "could not extract IP address from response"
		log.Println("Could not extract IP address from response")
		fmt.Println("------------------------------------------------------------")
		return result
	}

	externalIP := matches[1]

	// Validate IP address format
	if net.ParseIP(externalIP) == nil {
		result.Error = "invalid IP address extracted: " + externalIP
		log.Println("Invalid IP address extracted:", externalIP)
		fmt.Println("------------------------------------------------------------")
		return result
	}

	// Get local IP
	localIPs, err := net.LookupHost("localhost")
	if err != nil {
		result.Error = err.Error()
		log.Println("Error getting local IP:", err)
		fmt.Println("------------------------------------------------------------")
		return result
	}

	if len(localIPs) == 0 {
		result.Error = "no local IP addresses found"
		log.Println("No local IP addresses found")
		fmt.Println("------------------------------------------------------------")
		return result
	}

	if externalIP == localIPs[0] {
		result.Status = "Not using VPN or proxy."
		log.Println("Not using VPN or proxy.")
	} else {
		result.Status = "Using VPN or proxy."
		log.Println("Using VPN or proxy.")
	}

	fmt.Println("------------------------------------------------------------")
	return result
}
