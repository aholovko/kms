// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/kms/test/bdd

go 1.16

require (
	github.com/cucumber/godog v0.10.0
	github.com/fsouza/go-dockerclient v1.6.6
	github.com/hyperledger/aries-framework-go v0.1.7-0.20210429013345-a595aa0b19c4
	github.com/hyperledger/aries-framework-go/component/storageutil v0.0.0-20210427144858-06fb8b7d2d30
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20210422133815-2ef2d99cb692
	github.com/igor-pavlenko/httpsignatures-go v0.0.23
	github.com/rs/xid v1.2.1
	github.com/teserakt-io/golang-ed25519 v0.0.0-20210104091850-3888c087a4c8
	github.com/trustbloc/edge-core v0.1.7-0.20210429084532-c385778b4d8b
	github.com/trustbloc/edv v0.1.7-0.20210420141323-971877b36098
	github.com/trustbloc/hub-auth v0.1.7-0.20210429093653-7bf418749236 // indirect
	github.com/trustbloc/hub-auth/test/bdd v0.0.0-20210310171907-828969479ad5
	github.com/trustbloc/kms v0.0.0-00010101000000-000000000000
)

replace (
	github.com/trustbloc/kms => ../..
	// https://github.com/ory/dockertest/issues/208#issuecomment-686820414
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)
