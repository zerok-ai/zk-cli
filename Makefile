BINARY_NAME=zkcli
VERSION=0.0.1

build:
	GOARCH=amd64 GOOS=darwin  go build -o ./builds/${VERSION}/${BINARY_NAME}-darwin  -ldflags="-X 'root.version=${VERSION}'" main.go
	GOARCH=amd64 GOOS=linux   go build -o ./builds/${VERSION}/${BINARY_NAME}-linux   -ldflags="-X 'root.version=${VERSION}'" main.go
	GOARCH=amd64 GOOS=windows go build -o ./builds/${VERSION}/${BINARY_NAME}-windows -ldflags="-X 'root.version=${VERSION}'" main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -R ./builds