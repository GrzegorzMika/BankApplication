up:
	cd deployments && docker-compose up -d && cd ..

down:
	cd deployments && docker-compose down --remove-orphans && cd ..

migrate:
	cd db && soda migrate && cd ..

sqlc:
	sqlc generate

test:
	sudo go test -cover -v ./...

.PHONY: sqlc migrate up down test