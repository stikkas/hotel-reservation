build:
	@go build -o bin/api cmd/main.go

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test: 
	@go test -v ./...
