package gapi

import (
	"BankApplication/internal/db"
	"BankApplication/internal/util"
	"BankApplication/internal/worker"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}
