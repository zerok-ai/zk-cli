NAME=zkcli
VERSION=0.0.1

# Define the folder to delete
ARTIFACT_FOLDER := builds

delete-artifact-folder:
	@if [ -d "$(FOLDER_TO_DELETE)" ]; then \
		rm -rf "$(FOLDER_TO_DELETE)"; \
		echo "Folder deleted successfully."; \
	else \
		echo "Folder does not exist. No action taken."; \
	fi

FOLDER_TO_DELETE := $(ARTIFACT_FOLDER)
clean: delete-artifact-folder
	go clean

sync: clean
	go get -v ./...

build: sync
	go build -v -o $(NAME) main.go

run: build
	./${NAME}

artifact: build
	GOARCH=amd64 GOOS=darwin  go build -o ./${ARTIFACT_FOLDER}/${VERSION}/${NAME}-darwin  -ldflags="-X 'root.version=${VERSION}'" main.go
	GOARCH=amd64 GOOS=linux   go build -o ./${ARTIFACT_FOLDER}/${VERSION}/${NAME}-linux   -ldflags="-X 'root.version=${VERSION}'" main.go
	GOARCH=amd64 GOOS=windows go build -o ./${ARTIFACT_FOLDER}/${VERSION}/${NAME}-windows -ldflags="-X 'root.version=${VERSION}'" main.go
