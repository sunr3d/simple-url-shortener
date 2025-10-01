up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

clean:
	docker compose down -v

logs:
	docker compose logs -f app

test:
	go test -v ./...

migrate-up:
	docker compose exec -T db psql -U shortener_user -d shortener -f /migrations/001_init_up.sql

migrate-down:
	docker compose exec -T db psql -U shortener_user -d shortener -f /migrations/001_init_down.sql

fmt:
	go fmt ./...

build:
	go build -o simple-url-shortener cmd/main.go