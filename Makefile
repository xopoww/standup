.PHONY : gomod generate-proto standup-daemon standupctl

gomod:
	go mod tidy

generate-proto:
	./scripts/install_protoc.sh
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	./pkg/api/standup/generate

standup-daemon: gomod
	go build -o ./bin/standup-daemon ./cmd/standup-daemon

standupctl: gomod
	go build -o ./bin/standupctl ./cmd/standupctl