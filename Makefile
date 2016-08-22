CONFIG_DIR=/etc/bccp
DEFAULT_CONFIG_FILE=bccp.conf
PROJECT_NAME=bccp-server
SOURCES := $(shell find . -name '*.go' -print)

.PHONY: all clean install

all: $(PROJECT_NAME)

clean:
	$(RM) $(PROJECT_NAME)

fmt:
	go fmt ./...

install:
	mkdir -p $(CONFIG_DIR)
	cp $(DEFAULT_CONFIG_FILE) $(CONFIG_DIR)
	go install ./...

$(PROJECT_NAME): fmt $(SOURCES)
	go build

lint: ENABLE := vet vetshadow golint ineffassign gosimple
lint: EXCLUDE := 'comment.*exported' 'that|stutters' 'declaration|of|err|shadows|declaration'
lint: $(OUT_DIR) ##@lint Lint source code
	gometalinter --deadline=60s --disable-all $(addprefix --enable=,$(ENABLE)) $(subst |, ,$(addprefix --exclude=,$(EXCLUDE))) --sort=path --tests --vendor ./...