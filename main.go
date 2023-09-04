package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"BankApplication/internal/api"
	"BankApplication/internal/db"
	"BankApplication/internal/gapi"
	"BankApplication/internal/pb"
	"BankApplication/internal/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterBankApplicationHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalln("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalln("cannot create listener", err)
	}
	log.Println("HTTP gateway started ", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalln("cannot start HTTP gateway server: ", err)
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
