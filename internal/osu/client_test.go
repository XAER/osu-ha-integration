package osu_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/XAER/osu-ha-integration/internal/osu"
)

func TestGetAccessToken_CachesAndRefreshes(t *testing.T) {
	callCount := 0

	mockAuth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		res := map[string]interface{}{
			"access_token": "test-token",
			"expires_in":   10, // set a longer expiry
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}))
	defer mockAuth.Close()

	client := osu.NewClient("fake-id", "fake-secret")
	client.AuthURL = mockAuth.URL // override auth URL

	ctx := context.Background()

	// 1st call - should hit auth
	token1, err := client.GetAccessToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-token", token1)
	assert.Equal(t, 1, callCount)

	// 2nd call - should be cached
	token2, err := client.GetAccessToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, token1, token2)
	assert.Equal(t, 1, callCount)

	// Force expiry manually
	client.ExpiresAt = time.Now().Add(-time.Second)

	// 3rd call - should refresh
	token3, err := client.GetAccessToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, token3, "test-token")
	assert.Equal(t, 2, callCount)
}
