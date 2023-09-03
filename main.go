package main

import (
	"context"
	"log"
	"net"

	"BankApplication/internal/api"
	"BankApplication/internal/db"
	"BankApplication/internal/gapi"
	"BankApplication/internal/pb"
	"BankApplication/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalln("Cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankApplicationServer(grpcServer, server)

	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalln("cannot create listener", err)
	}
	log.Println("gRPC server listening on ", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("cannot start gRPC server: ", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalln("Cannot start server: ", err)
	}
}
