BINARY_NAME=zkcli

build:
	GOARCH=amd64 GOOS=darwin go build -o ./builds/${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ./builds/${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ./builds/${BINARY_NAME}-windows main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -R ./builds