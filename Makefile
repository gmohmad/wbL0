up:
	docker compose up --build

down:
	docker compose down

downv:
	docker compose down -v

gbuild:
	CGO_ENABLED=0 GOOS=linux go build -o ./builds/linux/ ./cmd/orders/main.go

run:
	go run ./cmd/orders/main.go

runb:
	./builds/linux/main

.PHONY: up down downv
