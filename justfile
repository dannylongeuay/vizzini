alias b := build
alias p := play
alias u := uci
alias s := serve
alias t := test
alias tl := test-long
alias l := lint
alias li := lichess
alias f := format

# Build binary
build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/vizzini ./src/...

# Run binary
play: build
    ./bin/vizzini play

uci: build
    ./bin/vizzini uci

serve: build
    CORS_PERMISSIVE=1 ./bin/vizzini serve

lichess: build
    ./bin/vizzini lichess

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
