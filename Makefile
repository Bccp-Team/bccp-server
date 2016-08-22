CONFIG_DIR=/etc/bccp
DEFAULT_CONFIG_FILE=bccp.conf
PROJECT_NAME=bccp-server
SOURCES := $(shell find . -name '*.go' -print)

.PHONY: all clean install

all: $(PROJECT_NAME)

clean:
	$(RM) $(PROJECT_NAME)

install: config
	mkdir -p $(CONFIG_DIR)
	cp $(DEFAULT_CONFIG_FILE) $(CONFIG_DIR)
	go install ./...

$(PROJECT_NAME): $(SOURCES)
	go build