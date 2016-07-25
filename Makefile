GOPATH_PREFIX=github.com
API_PKG=${GOPATH_PREFIX}/bccp/api
MAIN_PKG=${GOPATH_PREFIX}/bccp/main

GO=go
GO_FLAGS=install

all: ${API_PKG} ${MAIN_PKG}

Makefile: ;
.DEFAULT:
	${GO} ${GO_FLAGS} $@
