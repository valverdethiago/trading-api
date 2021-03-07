postgresql-start:
	docker-compose -f ./docker/docker-compose.yml up -d 

postgresql-stop:
	docker-compose -f ./docker/docker-compose.yml down 

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable" -verbose down

migrate-test-up:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade_test?sslmode=disable" -verbose up

migrate-test-down:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade_test?sslmode=disable" -verbose down

mockgen-query:
	mockgen -package mockdb -destination db/mock/query.go github.com/valverdethiago/trading-api/db/sqlc Querier

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: migrate-up migtrate-down postgresql-start postgresql-stop sqlc test server mockgen-query