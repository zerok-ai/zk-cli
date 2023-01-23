
# zk-cli sh

```sh

export GOPRIVATE=github.com/zerok-ai/*

# initialize the module
go mod init github.com/zerok-ai/zk-cli/zkctl


# create the repository structure 
mkdir -p commands
mkdir -p backend

touch commands/root.go
touch backend/logic.go
touch main.go

# add Cobra as a dependency
go get -u github.com/spf13/cobra@latest

```

All the commands should be in their respective files under the folder `.\commands` while the business logic should go in files in the folder `.\backend`.

```sh
touch commands/operator.go
touch commands/debug.go
```
