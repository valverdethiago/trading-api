migrateup:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:golang@localhost:6432/trade?sslmode=disable" -verbose down

.PHONY: migrateup migtratedown