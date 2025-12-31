package dsfhub

import (
	"errors"
	"strings"
)

// Config represents the configuration required for the DSF Client
type Config struct {
	// API Identifier
	DSFHUBToken string

	// API Key
	DSFHUBHost string

	// InsecureSSL
	InsecureSSL bool

	// Params including syncType, acknowledgeDeletionImpact, forceDelete
	Params map[string]interface{}
}

var validSyncTypes = []string{"SYNC_GW_BLOCKING", "SYNC_GW_NON_BLOCKING", "DO_NOT_SYNC_GW"}

var missingAPITokenMessage = "DSF HUB API Token must be provided"
var missingDSFHostMessage = "DSF HUB host/API endpoint must be provided"
var invalidSyncTypeMessage = "Invalid sync_type. Available values: " + strings.Join(validSyncTypes, ", ")

// Client configures and returns a fully initialized DSF Client
func (c *Config) Client() (interface{}, error) {
	// Check DSFToken
	if strings.TrimSpace(c.DSFHUBToken) == "" {
		return nil, errors.New(missingAPITokenMessage)
	}
	// Check DSFHost
	if strings.TrimSpace(c.DSFHUBHost) == "" {
		return nil, errors.New(missingDSFHostMessage)
	}
	// Check sync_type param
	if syncType, exists := c.Params["syncType"]; exists {
		if !isValidSyncType(syncType.(string)) {
			return nil, errors.New(invalidSyncTypeMessage)
		}
	}

	if acknowledgeDeletionImpact, exists := c.Params["acknowledgeDeletionImpact"]; exists {
		if _, ok := acknowledgeDeletionImpact.(string); !ok {
			return nil, errors.New("acknowledge_deletion_impact must be a string value of 'true' or 'false'")
		}
	}

	if forceDelete, exists := c.Params["forceDelete"]; exists {
		if _, ok := forceDelete.(string); !ok {
			return nil, errors.New("force_delete must be a string value of 'true' or 'false'")
		}
	}

	// Create client
	client := NewClient(c)

	// Verify client credentials
	gatewaysResponse, err := client.Verify()
	client.gateways = gatewaysResponse
	if err != nil {
		return nil, err
	}

	return client, nil
}

func isValidSyncType(sync_type string) bool {
	for _, valid_sync_type := range validSyncTypes {
		if sync_type == valid_sync_type {
			return true
		}
	}
	return false
}
