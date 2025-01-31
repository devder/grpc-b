package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/token"
	"github.com/devder/grpc-b/util"
	"github.com/devder/grpc-b/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
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

func newCtxWithBearerToken(t *testing.T, m token.Maker, username, role string, duration time.Duration) context.Context {
	accessToken, _, err := m.CreateToken(username, role, duration)
	require.NoError(t, err)
	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{bearerToken},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
