name = storageapi

main_path = ./cmd/storage-api/main.go

all: test build

build:
	go build -o ${name} ${main_path}

test:
	go test ./...

fclean:
	rm ${name}

re: fclean all

.PHONY: all test build fclean re