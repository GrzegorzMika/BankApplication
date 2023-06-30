up:
	cd deployments && docker-compose up -d && cd ..

down:
	cd deployments && docker-compose down --remove-orphans && cd ..

migrate:
	cd db && soda migrate && cd ..

sqlc:
	sqlc generate

.PHONY: sqlc migrate up down