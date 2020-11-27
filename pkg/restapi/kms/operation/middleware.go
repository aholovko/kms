/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package operation

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/piprate/json-gold/ld"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/zcapld"
)

// ZCAPLDMiddleware returns the ZCAPLD middleware that authorizes requests.
func (o *Operation) ZCAPLDMiddleware(h http.Handler) http.Handler {
	return &mwHandler{
		next:      h,
		zcaps:     o.authService,
		keys:      o.authService.KMS(),
		crpto:     o.authService.Crypto(),
		logger:    o.logger,
		routeFunc: (&muxNamer{}).GetName,
	}
}

type namer interface {
	GetName() string
}

type muxNamer struct {
}

func (m *muxNamer) GetName(r *http.Request) namer {
	return mux.CurrentRoute(r)
}

type mwHandler struct {
	next        http.Handler
	zcaps       zcapld.CapabilityResolver
	keys        kms.KeyManager
	crpto       crypto.Crypto
	logger      log.Logger
	ldDocLoader ld.DocumentLoader
	routeFunc   func(*http.Request) namer
}

func (h *mwHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling request: %s", r.URL.String())

	// this one is protected with OAuth2
	if h.routeFunc(r).GetName() == keystoresEndpoint {
		h.next.ServeHTTP(w, r)

		return
	}

	resource := keystoreLocation(r.Host, mux.Vars(r)[keystoreIDQueryParam])

	expectations := &zcapld.InvocationExpectations{
		Target:         resource,
		RootCapability: resource,
	}

	var err error

	expectations.Action, err = expectedAction(h.routeFunc(r))
	if err != nil {
		h.logger.Errorf("zcap middleware failed to determine the expected action: %s", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	// TODO make KeyResolver configurable
	// TODO make signature suites configurable
	zcapld.NewHTTPSigAuthHandler(
		&zcapld.HTTPSigAuthConfig{
			CapabilityResolver: h.zcaps,
			KeyResolver:        &zcapld.DIDKeyResolver{},
			VerifierOptions: []zcapld.VerificationOption{
				zcapld.WithSignatureSuites(
					ed25519signature2018.New(suite.WithVerifier(ed25519signature2018.NewPublicKeyVerifier())),
				),
				zcapld.WithLDDocumentLoaders(h.ldDocLoader),
			},
			Secrets:     &zcapld.AriesDIDKeySecrets{},
			ErrConsumer: h.logError,
			KMS:         h.keys,
			Crypto:      h.crpto,
		},
		expectations,
		h.next.ServeHTTP,
	).ServeHTTP(w, r)

	h.logger.Debugf("finished handling request: %s", r.URL.String())
}

func (h *mwHandler) logError(err error) {
	h.logger.Errorf("unauthorized capability invocation: %s", err.Error())
}

func expectedAction(n namer) (string, error) { // nolint:gocyclo // necessary complexity
	var (
		action string
		err    error
	)

	switch n.GetName() {
	case keysEndpoint:
		action = actionCreateKey
	case capabilityEndpoint:
		action = actionStoreCapability
	case exportEndpoint:
		action = actionExportKey
	case signEndpoint:
		action = actionSign
	case verifyEndpoint:
		action = actionVerify
	case encryptEndpoint:
		action = actionEncrypt
	case decryptEndpoint:
		action = actionDecrypt
	case computeMACEndpoint:
		action = actionComputeMac
	case verifyMACEndpoint:
		action = actionVerifyMAC
	case wrapEndpoint:
		action = actionWrap
	case unwrapEndpoint:
		action = actionUnwrap
	default:
		err = fmt.Errorf("unsupported endpoint: %s", n.GetName())
	}

	return action, err
}