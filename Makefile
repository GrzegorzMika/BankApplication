up:
	cd deployments && docker-compose up -d && cd ..

down:
	cd deployments && docker-compose down --remove-orphans && cd ..

migrate:
	cd db && soda migrate && cd ..

sqlc:
	sqlc generate

test:
	sudo go test -shuffle=on -cover -v ./...

test_with_coverage:
	sudo go test -shuffle=on -cover -v -coverprofile=coverage/coverage.out ./...

show_coverage:
	go tool cover -html=coverage/coverage.out

.PHONY: sqlc migrate up down test test_with_coverage show_coverage