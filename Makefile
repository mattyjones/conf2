export GOPATH=$(abspath .)
export YANGPATH=$(abspath etc)

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


libyangc2_CFLAGS = \
	-I$(abspath include)

libyangc2j_CFLAGS = \
	$(libyangc2_CFLAGS) \
	-I$(abspath drivers/java/include) \
	-I$(abspath include) \
	$(JDK_CFLAGS)

libyangc2j_LDFLAGS = \
	-L$(abspath pkg/$(GO_ARCH)_shared) -lyangc2 \
	$(JDK_LDFLAGS)

GO_ARCH = linux_amd64

PKGS = yang yang/browse

all : generate driver-java build test install

.PHONY: generate driver-java build test install
generate :
	go generate yang

build :
	CGO_CFLAGS="$(libyangc2_CFLAGS)" \
	  go build $(PKGS)

TEST='Test*'
Test% :
	$(MAKE) test TEST='$@'

test : src/yang/parser.go
	CGO_CFLAGS="$(libyangc2_CFLAGS)" \
	  go test -v $(PKGS) -run $(TEST)

install: libyangc2 libyangc2j;

JAVA_SRC = $(shell find drivers/java/src \( \
	-name '*.java' -a \
	-not -name '*Test.java' \) -type f)

JNI_SRCS = \
	org.conf2.yang.driver.Driver \
	org.conf2.yang.driver.DriverTestHarness \
	org.conf2.restconf.Service

clean :
	! test -d drivers/java/classes || rm -rf drivers/java/classes

driver-java :
	test -d drivers/java/classes || mkdir drivers/java/classes
	@javac -d drivers/java/classes $(JAVA_SRC)
	javah -cp drivers/java/classes -d drivers/java/include $(JNI_SRCS)
	jar -cf drivers/java/lib/yangc2-$(API_VER).jar -C drivers/java/classes .

JAVA_TEST_JARS = \
	$(wildcard drivers/java/lib/hamcrest-core-*.jar) \
	$(wildcard drivers/java/lib/hamcrest-library-*.jar) \
	$(wildcard drivers/java/lib/junit-*.jar)

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

debug :
	echo $(JAVA_TESTS)

driver-java-test :
	test -d drivers/java/test || mkdir drivers/java/test
	javac -d drivers/java/test -cp $(JAVA_TEST_CP) $(JAVA_TEST_SRC)
	LD_LIBRARY_PATH=$(JDK_LIBRARY_PATH):pkg/$(GO_ARCH)_shared \
	  CLASSPATH=drivers/java/src:drivers/java/test:$(JAVA_TEST_CP) \
	  java $(JAVA_TEST_RUNNER) $(JAVA_TESTS)

libyangc2j:
	cd pkg/$(GO_ARCH)_shared; \
	  ln -snf libyangc2j.a libyangc2j.so
	LD_LIBRARY_PATH=$(JDK_LIBRARY_PATH) \
	  CGO_CFLAGS="$($@_CFLAGS)" \
	  CGO_LDFLAGS="$($@_LDFLAGS)" \
	  go install -buildmode=c-shared $@

libyangc2:
	cd pkg/$(GO_ARCH)_shared; \
	  ln -snf libyangc2.a libyangc2.so
	CGO_CFLAGS="$($@_CFLAGS)" \
	  CGO_LDFLAGS="$($@_LDFLAGS)" \
	  go install -buildmode=c-shared $@
