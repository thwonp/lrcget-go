package lrclib

import (
	"net/http"
	"time"
	"crypto/tls"
	"net"
)

// Client represents an LRCLIB API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new LRCLIB client with enhanced security
func NewClient(baseURL string) *Client {
	// Configure TLS for security
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Configure transport with security settings
	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		DisableKeepAlives:   false,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		userAgent: "LRCGET v1.0.0 (https://github.com/tranxuanthang/lrcget)",
	}
}

// SetBaseURL sets the base URL for the client
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// GetBaseURL returns the current base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}
