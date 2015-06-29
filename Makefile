export GOPATH=$(abspath .)

all : generate build test

.PHONY: generate build test
generate :
	go generate yang

build :
	go build yang

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test : src/yang/parser.go
	go test -v yang -run $(TEST)
