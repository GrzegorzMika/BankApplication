up:
	cd deployments && docker-compose up -d && cd ..

down:
	cd deployments && docker-compose down --remove-orphans && cd ..

migrate:
	cd db && soda migrate && cd ..

sqlc:
	sqlc generate

generate_mocks:
	mockgen -package mockdb -destination internal/db/mocks/store.go BankApplication/internal/db Store

test:
	sudo go test -shuffle=on -cover -v ./...

test_with_coverage:
	sudo go test -shuffle=on -cover -v -coverprofile=coverage/coverage.out ./...

show_coverage:
	go tool cover -html=coverage/coverage.out

server:
	go run main.go

proto:
	rm -f internal/pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=internal/proto --go_out=internal/pb --go_opt=paths=source_relative \
        --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=internal/pb \
        --grpc-gateway_opt=paths=source_relative \
        --openapiv2_out=doc/swagger \
        --openapiv2_opt=allow_merge=true,merge_file_name=bank_application \
        internal/proto/*.proto

.PHONY: sqlc migrate up down test test_with_coverage show_coverage server generate_mocks proto