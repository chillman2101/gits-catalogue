server:
	go run cmd/app/main.go
build:
	go build cmd/app/main.go
seed:
	go run cmd/seed/main.go
migrate:
	go run cmd/migrate/main.go
swag:
	swag init -g cmd/app/main.go
download:
	go mod download
test:
	go test ./...
test-verbose:
	go test -v ./...
test-cover:
	go test ./... -cover