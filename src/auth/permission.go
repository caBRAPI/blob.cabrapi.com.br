package auth

import (
	"strings"

	"blob/src/config"
)

// Permission identifies an allowed operation.
type Permission string

const (
	PermRead   Permission = "read"
	PermWrite  Permission = "write"
	PermDelete Permission = "delete"
	PermAll    Permission = "all"
)

// APIKey is a token with an optional permission scope.
type APIKey struct {
	Token string
	Perms map[Permission]bool
}

// HasPerm reports whether the key allows the given operation.
func (k *APIKey) HasPerm(p Permission) bool {
	if k.Perms[PermAll] {
		return true
	}
	return k.Perms[p]
}

// parsePermissions converts a "read,write" list into a permission set.
func parsePermissions(raw string) map[Permission]bool {
	perms := make(map[Permission]bool)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "" {
			continue
		}
		perms[Permission(part)] = true
	}
	return perms
}

// defaultPermissions grants every operation.
func defaultPermissions() map[Permission]bool {
	return map[Permission]bool{PermAll: true}
}

// Keys returns the API keys configured for the process.
func Keys() []APIKey {
	var keys []APIKey
	if config.Env == nil {
		return keys
	}
	for _, entry := range strings.Split(config.Env.APIKeys, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		key := APIKey{Token: entry, Perms: defaultPermissions()}
		if idx := strings.IndexByte(entry, ':'); idx > 0 {
			key.Token = entry[:idx]
			key.Perms = parsePermissions(entry[idx+1:])
		}
		keys = append(keys, key)
	}
	return keys
}

// matchKey returns the API key matching the raw bearer token, or nil.
func matchKey(bearer string) *APIKey {
	for i := range Keys() {
		if Keys()[i].Token == bearer {
			return &Keys()[i]
		}
	}
	return nil
}

// masterToken is the legacy shared secret; it grants full access.
func masterToken() string {
	if config.Env == nil {
		return ""
	}
	return config.Env.TokenSecret
}

// Verify returns the API key when the bearer token is valid (either the master
// secret or a configured API key). Returns nil otherwise.
func Verify(bearerToken string) *APIKey {
	if bearerToken == "" {
		return nil
	}
	if masterToken() != "" && bearerToken == masterToken() {
		return &APIKey{Token: bearerToken, Perms: defaultPermissions()}
	}
	return matchKey(bearerToken)
}

// RequirePermission checks that the given key permits the operation.
// A nil key (public/unauthenticated) never passes.
func RequirePermission(key *APIKey, perms ...Permission) bool {
	if key == nil {
		return false
	}
	for _, p := range perms {
		if !key.HasPerm(p) {
			return false
		}
	}
	return true
}
