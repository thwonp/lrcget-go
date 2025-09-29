package lrclib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

// RequestChallenge requests a challenge from the LRCLIB API
func (c *Client) RequestChallenge(ctx context.Context) (*ChallengeResponse, error) {
	url := fmt.Sprintf("%s/api/request-challenge", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var challenge ChallengeResponse
	if err := decodeJSONResponse(resp, &challenge); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &challenge, nil
}

// SolveChallenge solves a proof-of-work challenge
func SolveChallenge(prefix, targetHex string) string {
	target, err := hex.DecodeString(targetHex)
	if err != nil {
		return ""
	}

	nonce := 0
	for {
		input := fmt.Sprintf("%s%d", prefix, nonce)
		hash := sha256.Sum256([]byte(input))

		if verifyNonce(hash[:], target) {
			break
		}
		nonce++
	}

	return fmt.Sprintf("%d", nonce)
}

// verifyNonce verifies if the hash meets the target requirement
func verifyNonce(result, target []byte) bool {
	if len(result) != len(target) {
		return false
	}

	for i := 0; i < len(result)-1; i++ {
		if result[i] > target[i] {
			return false
		} else if result[i] < target[i] {
			break
		}
	}

	return true
}
