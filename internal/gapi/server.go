package gapi

import (
	"fmt"

	"BankApplication/internal/db"
	"BankApplication/internal/pb"
	"BankApplication/internal/token"
	"BankApplication/internal/util"
)

// Server serves gRPC requests for banking application.
type Server struct {
	pb.UnimplementedBankApplicationServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}