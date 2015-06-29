export GOPATH=$(abspath .)

DESTDIR =
PREFIX = /usr/local/conf2
LIBDIR = $(PREFIX)/lib
INCLUDEDIR = $(PREFIX)/include
INSTALL = install

all : generate build test

.PHONY: build
generate :
	go generate yang

build:
	go build yang

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test : src/yang/parser.go
	go test -v yang -run $(TEST)

