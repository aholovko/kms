/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package keystore

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

// Keystore represents user's keystore with a list of associated keys and metadata.
type Keystore struct {
	ID                       string          `json:"id"`
	Controller               string          `json:"controller"`
	DelegateKeyID            string          `json:"delegateKeyID,omitempty"`
	RecipientKeyID           string          `json:"recipientKeyID,omitempty"`
	MACKeyID                 string          `json:"macKeyID,omitempty"`
	OperationalVaultID       string          `json:"operationalVaultID,omitempty"`
	OperationalEDVCapability json.RawMessage `json:"operationalEDVCapability,omitempty"`
	OperationalKeyIDs        []string        `json:"operationalKeyIDs,omitempty"`
	CreatedAt                *time.Time      `json:"createdAt"`
}

// Options configures Keystore during creation.
type Options struct {
	ID                 string
	Controller         string
	DelegateKeyType    kms.KeyType
	RecipientKeyType   kms.KeyType
	MACKeyType         kms.KeyType
	OperationalVaultID string
	CreatedAt          *time.Time
}

// Option configures Options.
type Option func(options *Options)

// WithID sets ID of Keystore.
func WithID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// WithController sets the controller of Keystore.
func WithController(c string) Option {
	return func(o *Options) {
		o.Controller = c
	}
}

// WithDelegateKeyType sets a type of the delegate key.
// Key is not created if type is not specified.
func WithDelegateKeyType(k kms.KeyType) Option {
	return func(o *Options) {
		o.DelegateKeyType = k
	}
}

// WithRecipientKeyType sets a type of the recipient key.
// Key is not created if type is not specified.
func WithRecipientKeyType(k kms.KeyType) Option {
	return func(o *Options) {
		o.RecipientKeyType = k
	}
}

// WithMACKeyType sets a type of the MAC key.
// Key is not created if type is not specified.
func WithMACKeyType(k kms.KeyType) Option {
	return func(o *Options) {
		o.MACKeyType = k
	}
}

// WithOperationalVaultID sets the ID of the vault on SDS server for storing operational keys.
func WithOperationalVaultID(id string) Option {
	return func(o *Options) {
		o.OperationalVaultID = id
	}
}

// WithCreatedAt sets the creation time of Keystore.
func WithCreatedAt(t *time.Time) Option {
	return func(o *Options) {
		o.CreatedAt = t
	}
}