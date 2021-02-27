postgresql-start:
	docker-compose -f ../docker/docker-compose.yml up -d 

postgresql-stop:
	docker-compose -f ../docker/docker-compose.yml down 

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: migrate-up migtrate-down postgresql-start postgresql-stop sqlc