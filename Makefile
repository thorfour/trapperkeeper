.PHONY: all clean push

EXE="server"
REPO="quay.io/thorfour/trapperkeeper"

all:  
	mkdir  -p ./bin/
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/$(EXE) ./cmd/server
	cp /etc/ssl/certs/ca-certificates.crt ./bin/
	cp ./build/Dockerfile ./bin/
	docker build ./bin/ -t $(REPO)
clean:
	rm -r ./bin
push:
	./build/docker_push.sh $(REPO)
