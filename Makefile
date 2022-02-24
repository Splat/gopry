DIRS=distribution bin bin/osx bin/linux_64 bin/win_64 bin/win_32 distribution
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
VERSION := $(shell git describe --tags)

.PHONY: test test-v cover fmt lint get build update make-dirs clean compile-osx compile-linux_64 compile-win_64 compile-win_32 compile-all osx-distribution linux-distribution win-distribution_64 win-distribution_32 distribution clean

default: build

test: build
	go test ./... -covermode=count -coverprofile=coverage/cover.out -timeout=5s

test-v:
	go test ./... -covermode=count -coverprofile=coverage/cover.out -timeout=5s -v

cover: test
	go tool cover -html=coverage/cover.out

fmt:
	goimports -w .

lint:
	golangci-lint run

get:
	go get

build: get
	go build -o gopry

update:
	go get -u -t
	go mod tidy

make-dirs:
	$(shell mkdir -p $(DIRS))

clean:
	rm -rf ./bin/
	rm -rf ./distribution/

compile-osx: make-dirs
	env GOOS=darwin GOARCH=amd64 go build -a -trimpath -ldflags="-X gopry/cmd.version=${VERSION} -s -w" -o ./bin/osx/jump

compile-linux_64: make-dirs
	env GOOS=linux GOARCH=amd64 go build -a -trimpath -ldflags="-X gopry/cmd.version=${VERSION} -s -w" -o ./bin/linux_64/jump

compile-win_64: make-dirs
	env GOOS=windows GOARCH=amd64 go build -a -trimpath -ldflags="-X gopry/cmd.version=${VERSION} -s -w" -o ./bin/win_64/jump.exe

compile-win_32: make-dirs
	env GOOS=windows GOARCH=386 go build -a -trimpath -ldflags="-X gopry/cmd.version=${VERSION} -s -w" -o ./bin/win_32/jump.exe

compile-all: compile-osx compile-linux_64 compile-win_64 compile-win_32

osx-distribution: compile-osx
	tar -cvzf ./distribution/jump-osx.tar.gz -C ./bin/osx/ gopry

linux-distribution: compile-linux_64
	gtar -cvzf ./distribution/jump-linux_64.tar.gz -C ./bin/linux_64/ gopry

win-distribution_64: compile-win_64
	zip -j ./distribution/jump-win_64.zip ./bin/win_64/gopry.exe

win-distribution_32: compile-win_32
	zip -j ./distribution/jump-win_32.zip ./bin/win_32/gopry.exe

distribution: compile-all osx-distribution linux-distribution win-distribution_64 win-distribution_32

release: clean distribution
