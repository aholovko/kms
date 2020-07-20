/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd

import (
	"errors"
	"net/http"
	"os"
	"testing"

	ariesmockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/hyperledger/aries-framework-go/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

const (
	testKeystoreID = "keystoreID"
)

type mockServer struct{}

func (s *mockServer) ListenAndServe(host string, handler http.Handler) error {
	return nil
}

func TestListenAndServe(t *testing.T) {
	var w HTTPServer
	err := w.ListenAndServe("wronghost", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "address wronghost: missing port in address")
}

func TestStartCmdContents(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	require.Equal(t, "start", startCmd.Use)
	require.Equal(t, "Start kms-rest", startCmd.Short)
	require.Equal(t, "Start kms-rest inside the hub-kms", startCmd.Long)

	checkFlagPropertiesCorrect(t, startCmd, hostURLFlagName, hostURLFlagShorthand, hostURLFlagUsage)
}

func TestStartCmdWithBlankArg(t *testing.T) {
	t.Run("test blank host url arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{"--" + hostURLFlagName, ""}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "host-url value is empty", err.Error())
	})
}

func TestStartCmdWithMissingArg(t *testing.T) {
	t.Run("test missing host url arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		err := startCmd.Execute()

		require.Error(t, err)
		require.Equal(t,
			"Neither host-url (command line flag) nor KMS_REST_HOST_URL (environment variable) have been set.",
			err.Error())
	})
}

func TestStartCmdWithBlankEnvVar(t *testing.T) {
	t.Run("test blank host env var", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		err := os.Setenv(hostURLEnvKey, "")
		require.NoError(t, err)

		err = startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "KMS_REST_HOST_URL value is empty", err.Error())
	})
}

func TestStartCmdValidArgs(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	args := []string{"--" + hostURLFlagName, "localhost:8080"}
	startCmd.SetArgs(args)

	err := startCmd.Execute()

	require.Nil(t, err)
}

func TestStartCmdValidArgsEnvVar(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	setEnvVars(t)

	defer unsetEnvVars(t)

	err := startCmd.Execute()
	require.NoError(t, err)
}

func TestCreateKeystoreProvider(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p, err := createKeystoreProvider()

		require.NotNil(t, p)
		require.NotNil(t, p.StorageProvider())
		require.NotNil(t, p.KMSCreator())
		require.NotNil(t, p.Crypto())
		require.NoError(t, err)
	})
}

func TestPrepareKMSCreator(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		kmsCreator := prepareKMSCreator(ariesmockstorage.NewMockStoreProvider())
		kms, err := kmsCreator(testKeystoreID)

		require.NotNil(t, kms)
		require.NoError(t, err)
	})

	t.Run("Error prepare master key reader", func(t *testing.T) {
		kmsCreator := prepareKMSCreator(
			&ariesmockstorage.MockStoreProvider{
				ErrOpenStoreHandle: errors.New("open store error")})

		kms, err := kmsCreator(testKeystoreID)

		require.Nil(t, kms)
		require.Error(t, err)
	})
}

func TestPrepareMasterKeyReader(t *testing.T) {
	t.Run("Error open store", func(t *testing.T) {
		reader, err := prepareMasterKeyReader(
			&ariesmockstorage.MockStoreProvider{
				ErrOpenStoreHandle: errors.New("open store error")})

		require.Nil(t, reader)
		require.Equal(t, errors.New("open store error"), err)
	})

	t.Run("Error retrieve master key from store", func(t *testing.T) {
		reader, err := prepareMasterKeyReader(
			&ariesmockstorage.MockStoreProvider{
				Store: &ariesmockstorage.MockStore{
					ErrGet: errors.New("get error")}})

		require.Nil(t, reader)
		require.Equal(t, errors.New("get error"), err)
	})

	t.Run("Error put newly generated master key into store", func(t *testing.T) {
		reader, err := prepareMasterKeyReader(
			&ariesmockstorage.MockStoreProvider{
				Store: &ariesmockstorage.MockStore{
					ErrGet: storage.ErrDataNotFound,
					ErrPut: errors.New("put error")}})

		require.Nil(t, reader)
		require.Equal(t, errors.New("put error"), err)
	})
}

func setEnvVars(t *testing.T) {
	err := os.Setenv(hostURLEnvKey, "localhost:8080")
	require.NoError(t, err)
}

func unsetEnvVars(t *testing.T) {
	err := os.Unsetenv(hostURLEnvKey)
	require.NoError(t, err)
}

func checkFlagPropertiesCorrect(t *testing.T, cmd *cobra.Command, flagName, flagShorthand, flagUsage string) {
	flag := cmd.Flag(flagName)

	require.NotNil(t, flag)
	require.Equal(t, flagName, flag.Name)
	require.Equal(t, flagShorthand, flag.Shorthand)
	require.Equal(t, flagUsage, flag.Usage)
	require.Equal(t, "", flag.Value.String())

	flagAnnotations := flag.Annotations
	require.Nil(t, flagAnnotations)
}
