gen_docs:
	swag init -g ./cmd/main/main.go -o ./docs;

build:
	go build -o bin/bin ./cmd/main/main.go && ./bin/bin;