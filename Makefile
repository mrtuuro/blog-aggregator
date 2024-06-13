
CONN="postgres://tuuro:@localhost:5432/blogator"

build:
	@go build -o bin/blog ./main.go

run: build
	@clear
	@./bin/blog


migration-up: 
	goose -dir ./sql/schema/ postgres postgres://tuuro:@localhost:5432/blogator up


migration-down: 
	goose -dir ./sql/schema/ postgres postgres://tuuro:@localhost:5432/blogator down

test:
	@go test ./...

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf ./bin

	
