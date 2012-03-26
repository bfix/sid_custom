
PWD := $(shell pwd)
BASE := $(shell dirname $(PWD))

ifneq ($(shell basename $(PWD)), src)
$(error Not in source directory - you are here: '${PWD}')
endif

ifneq ($(strip $(GOPATH)),)
GOPKG := $(shell echo ${GOPATH} | cut -f 1 -d ':')
GOPATH := ${GOPATH}:$(BASE)
GOINST := cp -vr $(BASE)/src/sid $(GOPKG)/bin
GOCLEAN := 	rm -f ${GOPKG}/bin/sid
else
GOPATH := $(BASE)
GOPKG := $(BASE)
GOINST := 
GOCLEAN := 
endif

install:
	go build -o sid
	$(GOINST)

test:
	go test

clean:
	GOPATH=${BASE} go clean
	rm -f ${BASE}/src/sid
	$(GOCLEAN)
