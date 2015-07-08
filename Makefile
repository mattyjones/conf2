export GOPATH=$(abspath .)

DESTDIR =
PREFIX = /usr/local/conf2
INCLUDEDIR = $(PREFIX)/include
LIBDIR = $(PREFIX)/lib
INSTALL = install

all : generate build test

.PHONY: generate build install test
generate :
	go generate yang

build :
	go build yang

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test : src/yang/parser.go
	go test -v yang -run $(TEST)

install:
	go install -buildmode=c-shared libyang
