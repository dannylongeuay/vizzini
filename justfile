alias b := build
alias r := run
alias t := test
alias tl := test-long
alias l := lint
alias f := format

# Build binary
build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/vizzini ./src/...

# Run binary
run: build
    ./bin/vizzini

# Run short tests
test:
    go test -v -count=1 -short ./src/...

# Run long tests
test-long:
    go test -v -count=1 ./src/...

# Lint all files
lint:
    golangci-lint run

# Format all files
format:
    goimports -w src
