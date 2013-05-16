
PWD := $(shell pwd)
SID := ${PWD}/../../SID
BASE := $(shell dirname $(PWD))

ifneq ($(shell basename $(PWD)), src)
$(error Not in source directory - you are here: '${PWD}')
endif

install:	fmt
	GOPATH=${SID}:${GOPATH} go build -o ../runtime/sid.app -v ./...

test:
	go test ./...

clean:
	go clean -i ./...
	
fmt:
	go fmt ./...
