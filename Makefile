BIN=bin/properties.bin
CMD=./cmd/properties
COVER=test.cover

GIT_HASH=`git rev-parse --short HEAD`
BUILD_DATE=`date +%FT%T%z`

LDFLAGS=-w -s -X main.GitHash=${GIT_HASH} -X main.BuildDate=${BUILD_DATE}

export CGO_ENABLED=0

.PHONY: clean build

build: vet
	go build -ldflags "${LDFLAGS}" -o "${BIN}" "${CMD}"

vet:
	go vet ./...

test:
	go test -race -count 1 -v -coverprofile="${COVER}" ./...

test-cover: test
	go tool cover -func="${COVER}"

lint:
	golangci-lint run

clean:
	[ -f "${BIN}" ] && rm "${BIN}"
	[ -f "${COVER}" ] && rm "${COVER}"
