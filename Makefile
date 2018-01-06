.PHONY: all clean

EXE="cmd"

all:  
	mkdir  -p ./bin/
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/$(EXE) ./cmd/*
	zip  -j ./bin/trapper.zip ./build/* ./bin/*
	zip -j ./bin/kick.zip ./build/kick/*
clean:
	rm -r ./bin
