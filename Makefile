export GOPATH=$(abspath .)

DESTDIR =
PREFIX = /usr/local/conf2
INCLUDEDIR = $(PREFIX)/include
LIBDIR = $(PREFIX)/lib
INSTALL = install

JDK_HOME = /usr/lib/jvm/java-8-oracle
JDK_OS = linux
JDK_ARCH = amd64
JDK_CFLAGS = -I$(JDK_HOME)/include -I$(JDK_HOME)/include/$(JDK_OS)
JDK_LDFLAGS = -L$(JDK_HOME)/jre/lib/$(JDK_ARCH) -ljava

libc2yang_CFLAGS = \
	-I$(abspath src)

libc2yangj_CFLAGS = \
	$(libc2yang_CFLAGS) \
	-I$(abspath drivers/java/include) \
	-I$(abspath pkg/$(GO_ARCH)_shared) \
	$(JDK_CFLAGS)

libc2yangj_LDFLAGS = \
	-L$(abspath pkg/$(GO_ARCH)_shared) -lc2yang \
	$(JDK_LDFLAGS)

GO_ARCH = linux_amd64

all : generate driver-java build test install

.PHONY: generate driver-java build test install
generate :
	go generate yang

build :
	CGO_CFLAGS="$(libc2yang_CFLAGS)" \
	  go build yang yang/browser yang/comm

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test : src/yang/parser.go
	CGO_CFLAGS="$(libc2yang_CFLAGS)" \
	  go test -v yang -run $(TEST)

install: ;
#install: libc2yang libc2yangj;

JAVA_SRC = $(shell find drivers/java/src \( \
	-name '*.java' -a \
	-not -name '*Test.java' \) -type f)

JNI_SRCS = \
	org.conf2.yang.Driver

driver-java :
	javac -d drivers/java/classes $(JAVA_SRC)
	javah -cp drivers/java/classes -d drivers/java/include $(JNI_SRCS)
	jar -cf c2yang.jar -C drivers/java/classes .


libc2yang libc2yangj:
	CGO_CFLAGS="$($@_CFLAGS)" \
	  CGO_LDFLAGS="$($@_LDFLAGS)" \
	  go install -buildmode=c-shared $@
