package utils

import "time"

// TestResults represents the complete results of all tests
type TestResults struct {
	HTTPTests  []HTTPTest  `json:"http_tests,omitempty"`
	SpeedTests []SpeedTest `json:"speed_tests,omitempty"`
	VPNTest    VPNTest     `json:"vpn_test,omitempty"`
	PingTest   PingTest    `json:"ping_test,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

// HTTPTest represents the result of an HTTP test
type HTTPTest struct {
	URL            string `json:"url"`
	Status         string `json:"status"`
	Proto          string `json:"proto,omitempty"`
	TLSVersion     string `json:"tls_version,omitempty"`
	CipherSuite    string `json:"cipher_suite,omitempty"`
	ServerName     string `json:"server_name,omitempty"`
	ResponseLength int    `json:"response_length,omitempty"`
	Error          string `json:"error,omitempty"`
}

// SpeedTest represents the result of a speed test
type SpeedTest struct {
	URL           string        `json:"url"`
	DownloadMbps  float64       `json:"download_mbps"`
	ElapsedTime   time.Duration `json:"elapsed_time"`
	BytesReceived int           `json:"bytes_received"`
	Error         string        `json:"error,omitempty"`
}

// VPNTest represents the result of a VPN detection test
type VPNTest struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// PingTest represents the result of a ping test
type PingTest struct {
	URL         string  `json:"url,omitempty"`
	Transmitted int     `json:"transmitted_packets,omitempty"`
	Received    int     `json:"received_packets,omitempty"`
	Loss        float64 `json:"loss_packets,omitempty"`
	Error       string  `json:"error,omitempty"`
}

// Tests is kept for backward compatibility with existing data.json
type Tests struct {
	VPNTest  VPNTest  `json:"vpn_test"`
	PingTest PingTest `json:"ping_test"`
}
