package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/palSagnik/go-YTFetch.git/backend/config"
	"github.com/palSagnik/go-YTFetch.git/backend/database"
	"github.com/palSagnik/go-YTFetch.git/backend/utils"
)

type APIKeyManager struct {
	keys       []string
	quotaUsage map[string]int
	resetTime  map[string]time.Time
	mu         sync.Mutex
}

const (
	maxDailyQuota      = 10000 // YouTube's default daily quota limit
	quotaCostPerSearch = 100   // Cost for each search operation
)

// creates a new API key manager with the given keys
func NewAPIKeyManager(keys []string) (*APIKeyManager, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no API keys provided")
	}

	return &APIKeyManager{
		keys:       keys,
		quotaUsage: make(map[string]int),
		resetTime:  make(map[string]time.Time),
	}, nil
}

// returning the key with the lowest quota usage
func (m *APIKeyManager) GetNextKey() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var bestKey string
	minQuota := maxDailyQuota + 1

	// Find key with lowest quota
	for _, key := range m.keys {

		// Check if quota should be reset
		if resetTime, exists := m.resetTime[key]; exists {
			if time.Now().After(resetTime) {
				m.quotaUsage[key] = 0
				delete(m.resetTime, key)
			}
		}

		quota := m.quotaUsage[key]
		if quota < minQuota {
			minQuota = quota
			bestKey = key
		}
	}

	// Check if we found a usable key
	if minQuota >= maxDailyQuota {
		return "", fmt.Errorf("all API keys have exceeded their quota")
	}

	// Update quota usage for selected key
	m.quotaUsage[bestKey] += quotaCostPerSearch

	// Set reset time if this is first usage
	if _, exists := m.resetTime[bestKey]; !exists {
		m.resetTime[bestKey] = time.Now().Add(24 * time.Hour)
	}

	return bestKey, nil
}

// marking the keys whose quotas have exceeded
func (m *APIKeyManager) MarkQuotaExceeded(key string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.quotaUsage[key] = maxDailyQuota
    m.resetTime[key] = time.Now().Add(24 * time.Hour)
}

func FetchingCronJob() {
	
	// predefined search query strings
	// based on google youtube trends in india
	queryKeywords := []string{
		"cricket",
		"politics",
		"entertainment",
		"vlogs",
		"music",
	}

	// apikeys for usage
	keys := []string{
		config.YOUTUBE_APIKEY_1,
		config.YOUTUBE_APIKEY_2,
		config.YOUTUBE_APIKEY_3,
	}

	// initialising new apikey manager
	keyManager, err := NewAPIKeyManager(keys)
	if err != nil {
		fmt.Printf("error in creating ApiKeyManager: %s", err.Error())
	}

	for {
		fmt.Println("### STARTING CRONJOB ###")
		
		// added this to prevent apikey quota exhaustion
		// comment this line while starting the server
		time.Sleep(500 * time.Hour)

		apikey, err := keyManager.GetNextKey()
		if err != nil {
			fmt.Printf("error in getting api key: %s", err.Error())
		}

		for _, query := range queryKeywords {
			
			var ctx context.Context

			videos, err := utils.SearchYoutubevideos(apikey, query, config.DEFAULT_PUBLISHED_AFTER, int64(config.DEFAULT_VIDEO_FETCH_LIMIT))
			if err != nil {
				fmt.Printf("error in fetching videos of query(%s): %s", query, err.Error())
			}

			err = database.InsertVideoDetailsCron(ctx, videos)
			if err != nil {
				fmt.Printf("error in inserting videos of query(%s): %s", query, err.Error())
			}
		}

		//time.Sleep(time.Hour * 1)
	}

}
