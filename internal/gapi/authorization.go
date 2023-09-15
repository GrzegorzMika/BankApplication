package gapi

import (
	"context"
	"fmt"
	"strings"

	"BankApplication/internal/token"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata from incoming context")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fileds := strings.Fields(authHeader)
	if len(fileds) != 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}
	authType := strings.ToLower(fileds[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}
	accessToken := fileds[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %s", err)
	}
	return payload, nil
}
