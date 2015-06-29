all : src/yang/parser.go build test

build : $(wildcard yang/*.go)
	go build yang

lib : $(wildcard yang/*.go)
	go build -buildmode=c-archive yang

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test : src/yang/parser.go
	go test -v yang -run $(TEST)


src/yang/parser.go : src/yang/parser.y
	go tool yacc -o $@.tmp $^ && \
		mv $@.tmp $@
