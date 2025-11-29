package modules

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ehsanghaffar/ultimate-internet-test/config"
	"github.com/ehsanghaffar/ultimate-internet-test/utils"
	"github.com/go-ping/ping"
)

// PingCheck performs a ping test on the given domain and returns the result.
// It accepts a config parameter for ping configuration and returns a PingTest result with any errors.
// The function sends ping packets as configured and collects statistics about packet loss and timing.
//
// Parameters:
//   - domain: The domain or IP address to ping
//   - cfg: Configuration containing ping count and other settings
//
// Returns:
//   - *PingTest: Pointer to PingTest struct containing ping statistics and any errors
//
// Example:
//
//	cfg := config.New()
//	result := PingCheck("example.com", cfg)
//	if result.Error == "" {
//	    log.Printf("Packets: %d sent, %d received, %.2f%% loss\n",
//	        result.Transmitted, result.Received, result.Loss)
//	}
func PingCheck(domain string, cfg *config.Config) *utils.PingTest {
	result := &utils.PingTest{
		URL: domain,
	}

	pinger, err := ping.NewPinger(domain)
	if err != nil {
		result.Error = err.Error()
		log.Printf("Failed to create pinger for %s: %v\n", domain, err)
		fmt.Println("------------------------------------------------------------")
		return result
	}

	// Set ping count from config
	pinger.Count = cfg.PingCount

	// Listen for Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		log.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		log.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		// Update result with final statistics
		result.Transmitted = stats.PacketsSent
		result.Received = stats.PacketsRecv
		result.Loss = stats.PacketLoss
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	if err := pinger.Run(); err != nil {
		result.Error = err.Error()
		log.Printf("Ping check failed for %s: %v\n", domain, err)
		fmt.Println("------------------------------------------------------------")
		return result
	}

	fmt.Println("------------------------------------------------------------")
	return result
}
