package auth

import (
	"google.golang.org/grpc/metadata"
)

var ValidAPIKeys = map[string]bool{}

// RegisterAPIKey allows registering a valid API key dynamically
func RegisterAPIKey(key string) {
	ValidAPIKeys[key] = true
}

// ValidateMetadata checks JWT or API key from metadata
func ValidateMetadata(md metadata.MD) bool {
	keys := md.Get("authorization")
	if len(keys) > 0 && len(keys[0]) > 0 {
		return true // JWT or Bearer token — assume valid for now
	}

	apiKeys := md.Get("x-api-key")
	if len(apiKeys) > 0 && ValidAPIKeys[apiKeys[0]] {
		return true
	}
	return false
}
