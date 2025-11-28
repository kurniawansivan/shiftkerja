DB_URL=postgresql://postgres:password123@localhost:5432/shiftkerja?sslmode=disable

postgres:
	docker compose up -d

createdb:
	docker compose exec postgres createdb --username=postgres --owner=postgres shiftkerja

dropdb:
	docker compose exec postgres dropdb shiftkerja

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown