NAME=zkctl
VERSION?=0.0.1

#CLOUD_ADDRESS=devcloud01.getanton.com
CLOUD_ADDRESS=sandbox.zerok.dev

# Define the folder to delete
ARTIFACT_FOLDER_NAME ?= builds

delete-artifact-folder:
	@if [ -d "$(FOLDER_TO_DELETE)" ]; then \
		rm -rf "$(FOLDER_TO_DELETE)"; \
		echo "Folder deleted successfully."; \
	else \
		echo "Folder does not exist. No action taken."; \
	fi

FOLDER_TO_DELETE := $(ARTIFACT_FOLDER_NAME)
clean: delete-artifact-folder
	go clean

sync: clean
	go get -v ./...

build: sync
	echo "building for version = $(VERSION)"
	go build -o $(NAME) -ldflags="-X 'zkctl/cmd.BinaryVersion=$(VERSION) -X 'zkctl/cmd.prodCloudAddress=$(CLOUD_ADDRESS)'" main.go

run: build
	./$(NAME)

artifact:
	echo "building for version = $(VERSION)"
	GOARCH=amd64 GOOS=darwin  go build -o ./$(ARTIFACT_FOLDER_NAME)/$(VERSION)/$(NAME)-$(VERSION)-darwin  -ldflags="-X 'zkctl/cmd.BinaryVersion=$(VERSION)' -X 'zkctl/cmd.prodCloudAddress=$(CLOUD_ADDRESS)'" main.go
	GOARCH=amd64 GOOS=linux   go build -o ./$(ARTIFACT_FOLDER_NAME)/$(VERSION)/$(NAME)-$(VERSION)-linux   -ldflags="-X 'zkctl/cmd.BinaryVersion=$(VERSION)' -X 'zkctl/cmd.prodCloudAddress=$(CLOUD_ADDRESS)'" main.go
	#GOARCH=amd64 GOOS=windows go build -o ./$(ARTIFACT_FOLDER_NAME)/$(VERSION)/$(NAME)-$(VERSION)-windows -ldflags="-X 'zkctl/cmd.BinaryVersion=$(VERSION)' -X 'zkctl/cmd.prodCloudAddress=$(CLOUD_ADDRESS)'" main.go

delete:
	go run main.go delete -y

delete-artifact: artifact
	./$(ARTIFACT_FOLDER_NAME)/$(VERSION)/$(NAME)-darwin delete -y

run-prod:
	ZK_CLOUD_ADDRESS=$(CLOUD_ADDRESS)
	ZK_API_KEY=px-api-e0593597-de51-44cd-bc72-6cbdb881b2be
	echo "-"
	go run main.go install -y --apikey $(ZK_API_KEY)

run-dev:
	ZK_CLOUD_ADDRESS=devcloud01.getanton.com
	ZK_CLIENT_VERSION=0.1.0-alpha
	ZK_API_KEY=px-api-e0593597-de51-44cd-bc72-6cbdb881b2be
	echo "-"
	go run main.go install -y --apikey $(ZK_API_KEY) --dev --zkVersion=zk-scenario-manager=$(ZK_CLIENT_VERSION),zk-axon=$(ZK_CLIENT_VERSION),zk-daemonset=$(ZK_CLIENT_VERSION),zk-gpt=$(ZK_CLIENT_VERSION),zk-wsp-client=$(ZK_CLIENT_VERSION),zk-operator=$(ZK_CLIENT_VERSION),zk-app-init-containers=$(ZK_CLIENT_VERSION)

run-dev-artifact: artifact
	ZK_CLOUD_ADDRESS=devcloud01.getanton.com
	ZK_CLIENT_VERSION=0.1.0-alpha
	ZK_API_KEY=px-api-e0593597-de51-44cd-bc72-6cbdb881b2be
	echo "-"
	./$(ARTIFACT_FOLDER_NAME)/$(VERSION)/$(NAME)-darwin install -y --dev --apikey $(ZK_API_KEY) --zkVersion=zk-scenario-manager=$(ZK_CLIENT_VERSION),zk-axon=$(ZK_CLIENT_VERSION),zk-daemonset=$(ZK_CLIENT_VERSION),zk-gpt=$(ZK_CLIENT_VERSION),zk-wsp-client=$(ZK_CLIENT_VERSION),zk-operator=$(ZK_CLIENT_VERSION),zk-app-init-containers=$(ZK_CLIENT_VERSION)

ci-cd-artifact: artifact

ci-cd-artifact-version: delete-artifact-folder
	mkdir -p $(ARTIFACT_FOLDER_NAME)
	printf "${VERSION}" > $(ARTIFACT_FOLDER_NAME)/version.txt

ci-cd-helm-version: delete-artifact-folder
	mkdir -p $(ARTIFACT_FOLDER_NAME)
	printf "${VERSION}" > $(ARTIFACT_FOLDER_NAME)/helmversion.txt

ci-cd-ebpf-helm-version: delete-artifact-folder
	mkdir -p $(ARTIFACT_FOLDER_NAME)
	printf "${VERSION}" > $(ARTIFACT_FOLDER_NAME)/ebpfversion.txt

ci-cd-artifact-install: delete-artifact-folder
	mkdir -p $(ARTIFACT_FOLDER_NAME)
	cp ./install.sh $(ARTIFACT_FOLDER_NAME)/install.sh