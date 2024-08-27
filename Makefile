build: 
	@go build -o ./bin/go-market ./cmd/main/main.go

run: build
	@./bin/go-market