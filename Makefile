CURDIR:=$(shell pwd)
all:
	export GOPATH=$(CURDIR):$(GOPATH) && \
		go install main

clean:
	rm -rfv $(CURDIR)/pkg $(CURDIR)/bin/main
