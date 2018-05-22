.PHONY: all clean push docker plugin

EXE="server"
REPO="quay.io/thorfour/trapperkeeper"
PLUG="pick"

docker:
	mkdir -p ./bin/docker
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/docker/$(EXE) ./cmd/server/
	cp /etc/ssl/certs/ca-certificates.crt ./bin/docker
	cp ./build/Dockerfile ./bin/docker
	docker build ./bin/docker -t $(REPO)
clean:
	rm -r ./bin
push:
	./build/docker_push.sh $(REPO)
plugin: 
	mkdir -p ./bin/plugin
	CGO_ENABLED=0 GOOS=linux go build -buildmode=plugin -o ./bin/plugin/$(PLUG) ./cmd/plugin/
all: docker plugin
