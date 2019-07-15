.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/todos internal/lambda/todos/main.go
	chmod +x bin/todos

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
