build:
	@go build -o bin/blog ./main.go

run: build
	@clear
	@./bin/blog

test:
	@go test ./...

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf ./bin

	
