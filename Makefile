# unix only atm

ifeq ($(GOOS),windows)
EXT = ".exe"
else
SRC=$(shell find . -type d \( -path ./vendor -o -path ./testdata \) -prune -o -name '*.go' -print)
endif


.PHONY: build
build: kooky

.PHONY: all
all: kooky


kooky: ${SRC}
	@env GOWORK=off go build -trimpath -ldflags '-w -s' -o kooky${EXT} ./cmd/kooky


.PHONY: test
test:
	@env GOWORK=off go test -count=1 -timeout=30s ./... | grep -v '^? '

.PHONY: clean
clean:
	@rm -f -- kooky kooky.exe kooky.test
