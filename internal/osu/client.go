package osu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/XAER/osu-ha-integration/internal/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewClient(token string) *Client {
	return &Client{httpClient: &http.Client{Timeout: 5 * time.Second}, baseURL: "https://osu.ppy.sh/api/v2", token: token}
}

func (c *Client) GetUser(ctx context.Context, username string) (*domain.OsuUser, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", c.baseURL, username), nil)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	user := &domain.OsuUser{
		Username:    raw["username"].(string),
		GlobalRank:  int(raw["global_rank"].(float64)),
		CountryRank: int(raw["country_rank"].(float64)),
		PP:          raw["pp"].(float64),
		Accuracy:    raw["accuracy"].(float64),
		PlayCount:   int(raw["play_count"].(float64)),
	}

	return user, nil
}
