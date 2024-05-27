up:
	docker compose up --build

down:
	docker compose down

downv:
	docker compose down -v

gbuild:
	CGO_ENABLED=0 GOOS=linux go build -o ./builds/linux/ ./cmd/orders/main.go
	CGO_ENABLED=0 GOOS=darwin go build -o ./builds/macos/ ./cmd/orders/main.go
	CGO_ENABLED=0 GOOS=windows go build -o ./builds/windows/ ./cmd/orders/main.go

grun:
	go run ./cmd/orders/main.go

grub:
	./builds/linux/main

.PHONY: up down
