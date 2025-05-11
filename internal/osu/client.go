package osu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/XAER/osu-ha-integration/internal/domain"
)

type Client struct {
	httpClient   *http.Client
	BaseURL      string
	AuthURL      string
	clientID     string
	clientSecret string

	accessToken string
	ExpiresAt   time.Time
	mu          sync.Mutex
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		httpClient:   &http.Client{Timeout: 5 * time.Second},
		BaseURL:      "https://osu.ppy.sh/api/v2",
		AuthURL:      "https://osu.ppy.sh/oauth/token",
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.ExpiresAt) && c.accessToken != "" {
		return c.accessToken, nil // still valid
	}

	payload := map[string]interface{}{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"grant_type":    "client_credentials",
		"scope":         "public",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", c.AuthURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.Status)
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	c.accessToken = res.AccessToken
	c.ExpiresAt = time.Now().Add(time.Duration(res.ExpiresIn-5) * time.Second)

	return c.accessToken, nil
}

func (c *Client) GetUser(ctx context.Context, username string) (*domain.OsuUser, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", c.BaseURL, username), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("osu! API returned %s", resp.Status)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &domain.OsuUser{
		Username:    raw["username"].(string),
		GlobalRank:  int(raw["statistics"].(map[string]interface{})["global_rank"].(float64)),
		CountryRank: int(raw["statistics"].(map[string]interface{})["country_rank"].(float64)),
		PP:          raw["statistics"].(map[string]interface{})["pp"].(float64),
		Accuracy:    raw["statistics"].(map[string]interface{})["hit_accuracy"].(float64),
		PlayCount:   int(raw["statistics"].(map[string]interface{})["play_count"].(float64)),
	}, nil
}
