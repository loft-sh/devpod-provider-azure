//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package policy

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// BearerTokenOptions configures the bearer token policy's behavior.
type BearerTokenOptions struct {
	// Scopes contains the list of permission scopes required for the token.
	Scopes []string
}

// RegistrationOptions configures the registration policy's behavior.
// All zero-value fields will be initialized with their default values.
type RegistrationOptions struct {
	policy.ClientOptions

	// MaxAttempts is the total number of times to attempt automatic registration
	// in the event that an attempt fails.
	// The default value is 3.
	// Set to a value less than zero to disable the policy.
	MaxAttempts int

	// PollingDelay is the amount of time to sleep between polling intervals.
	// The default value is 15 seconds.
	// A value less than zero means no delay between polling intervals (not recommended).
	PollingDelay time.Duration

	// PollingDuration is the amount of time to wait before abandoning polling.
	// The default valule is 5 minutes.
	// NOTE: Setting this to a small value might cause the policy to prematurely fail.
	PollingDuration time.Duration
}

// ClientOptions contains configuration settings for a client's pipeline.
type ClientOptions struct {
	policy.ClientOptions

	// DisableRPRegistration disables the auto-RP registration policy. Defaults to false.
	DisableRPRegistration bool
}

// Clone return a deep copy of the current options.
func (o *ClientOptions) Clone() *ClientOptions {
	if o == nil {
		return nil
	}
	copiedOptions := *o
	copiedOptions.Cloud.Services = copyMap(copiedOptions.Cloud.Services)
	copiedOptions.Logging.AllowedHeaders = copyArray(copiedOptions.Logging.AllowedHeaders)
	copiedOptions.Logging.AllowedQueryParams = copyArray(copiedOptions.Logging.AllowedQueryParams)
	copiedOptions.Retry.StatusCodes = copyArray(copiedOptions.Retry.StatusCodes)
	copiedOptions.PerRetryPolicies = copyArray(copiedOptions.PerRetryPolicies)
	copiedOptions.PerCallPolicies = copyArray(copiedOptions.PerCallPolicies)
	return &copiedOptions
}

// copyMap return a new map with all the key value pair in the src map
func copyMap[K comparable, V any](src map[K]V) map[K]V {
	if src == nil {
		return nil
	}
	copiedMap := make(map[K]V)
	for k, v := range src {
		copiedMap[k] = v
	}
	return copiedMap
}

// copyMap return a new array with all the elements in the src array
func copyArray[T any](src []T) []T {
	if src == nil {
		return nil
	}
	copiedArray := make([]T, len(src))
	copy(copiedArray, src)
	return copiedArray
}
