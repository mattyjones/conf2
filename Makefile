export GOPATH=$(abspath .)
export YANGPATH=$(abspath etc/yang)

API_VER = 0.1
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
JDK_LIBRARY_PATH = $(JDK_HOME)/jre/lib/$(JDK_ARCH)/server:$(JDK_HOME)/jre/lib/$(JDK_ARCH)


libconf2_CFLAGS = \
	-I$(abspath include)

libconf2j_CFLAGS = \
	$(libconf2_CFLAGS) \
	-I$(abspath drivers/java/include) \
	-I$(abspath include) \
	$(JDK_CFLAGS)

libconf2j_LDFLAGS = \
	-L$(abspath pkg/$(GO_ARCH)_shared) -lconf2 \
	$(JDK_LDFLAGS)

GO_ARCH = linux_amd64

PKGS = \
	conf2 \
	schema \
	schema/yang \
	data \
	comm \
	process \
	process/yapl \
	restconf \
	app

all : generate driver-java build test install

.PHONY: generate driver-java build test install
generate : generate-yang-parser generate-yapl-parser

generate-yang-parser :
	go generate schema/yang

generate-yapl-parser :
	go generate process/yapl

build :
	CGO_CFLAGS="$(libconf22_CFLAGS)" \
	  go build $(PKGS)

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test :
	CGO_CFLAGS="$(libconf2_CFLAGS)" \
	  go test -v $(PKGS) -run $(TEST)

go-install :
	  go install $(PKGS)

deps :;

install: go-install libconf2 libconf2j;

JAVA_SRC = $(shell find drivers/java/src \( \
	-name '*.java' -a \
	-not -name '*Test.java' \) -type f)

JNI_SRCS = \
	org.conf2.schema.driver.Driver \
	org.conf2.schema.yang.ModuleLoader \
	org.conf2.schema.driver.DriverTestHarness \
	org.conf2.restconf.Service

clean :
	! test -d pkg || rm -rf pkg
	! test -d drivers/java/classes || rm -rf drivers/java/classes

fmt : 
	go fmt $(PKGS)

driver-java :
	test -d drivers/java/classes || mkdir drivers/java/classes
	@javac -d drivers/java/classes $(JAVA_SRC)
	javah -cp drivers/java/classes -d drivers/java/include $(JNI_SRCS)
	jar -cf drivers/java/lib/conf2-$(API_VER).jar -C drivers/java/classes .

JAVA_TEST_JARS = \
	$(wildcard drivers/java/lib/hamcrest-core-*.jar) \
	$(wildcard drivers/java/lib/hamcrest-library-*.jar) \
	$(wildcard drivers/java/lib/junit-*.jar)

remove-debug-prints :
	grep -rl conf2.Debug src | xargs sed -i -e '/conf2.Debug.Printf/d'

EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
JAVA_TEST_CP = drivers/java/classes:$(subst $(SPACE),:,$(JAVA_TEST_JARS))
JAVA_TEST_RUNNER = org.junit.runner.JUnitCore

JAVA_TESTS = Test

JAVA_TEST_SRC = \
	$(shell find drivers/java/src -name '*Test.java')

JAVA_TEST_SRC_BASE = \
	$(shell find drivers/java/src -name '*Test.java' -printf '%P ')

JAVA_TESTS = \
	$(subst /,.,$(JAVA_TEST_SRC_BASE:.java=))

driver-java-test :
	test -d drivers/java/test || mkdir drivers/java/test
	javac -d drivers/java/test -cp $(JAVA_TEST_CP) $(JAVA_TEST_SRC)
	LD_LIBRARY_PATH=$(JDK_LIBRARY_PATH):pkg/$(GO_ARCH)_shared \
	  CLASSPATH=drivers/java/src:drivers/java/test:$(JAVA_TEST_CP) \
	  java $(JAVA_TEST_RUNNER) $(JAVA_TESTS)

libconf2j:
	cd pkg/$(GO_ARCH)_shared; \
	  ln -snf libconf2j.a libconf2j.so
	LD_LIBRARY_PATH=$(JDK_LIBRARY_PATH) \
	  CGO_CFLAGS="$($@_CFLAGS)" \
	  CGO_LDFLAGS="$($@_LDFLAGS)" \
	  go install -buildmode=c-shared $@

libconf2:
	cd pkg/$(GO_ARCH)_shared; \
	  ln -snf libconf2.a libconf2.so
	CGO_CFLAGS="$($@_CFLAGS)" \
	  CGO_LDFLAGS="$($@_LDFLAGS)" \
	  go install -buildmode=c-shared $@
