install :
	GOPATH=$(abspath .):$(abspath ../..) \
	  go install todo

clean :
	! test -d pkg || rm -rf pkg
