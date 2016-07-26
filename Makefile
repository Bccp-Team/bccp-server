GOPATH_PREFIX=github.com
API_PKG=${GOPATH_PREFIX}/bccp/api
MAIN_PKG=${GOPATH_PREFIX}/bccp/bccp

CONFIG_DIR=/etc/bccp
DEFAULT_CONFIG_FILE=bccp.conf

GO=go
GO_FLAGS=install

all: ${API_PKG} ${MAIN_PKG}

install: config
config: ${CONFIG_DIR}
	cp ${DEFAULT_CONFIG_FILE} ${CONFIG_DIR}

${CONFIG_DIR}:
	mkdir -p ${CONFIG_DIR}


Makefile: ;
.DEFAULT:
	${GO} ${GO_FLAGS} $@
