# looper

looper is a CLI tool that allows the creation of loops.

## Install

### CLI

Setting target destination:

`curl -s https://raw.githubusercontent.com/thalesfsp/looper/main/resources/install.sh | BIN_DIR=ABSOLUTE_DIR_PATH sh`

Setting version:

`curl -s https://raw.githubusercontent.com/thalesfsp/looper/main/resources/install.sh | VERSION=v{M.M.P} sh`

Example:

`curl -s https://raw.githubusercontent.com/thalesfsp/looper/main/resources/install.sh | BIN_DIR=/usr/local/bin VERSION=v0.0.1 sh`

### Programmatically

Install dependency:

`go get -u github.com/thalesfsp/looper`

## Usage

### CLI

`looper l v --help` 

### Programmatically

See [`looper/looper_test.go`](looper/looper_test.go)

### Documentation

Run `$ make doc` or check out [online](https://pkg.go.dev/github.com/thalesfsp/looper).

## Development

Check out [CONTRIBUTION](CONTRIBUTION.md).

### Release

1. Update [CHANGELOG](CHANGELOG.md) accordingly.
2. Once changes from MR are merged.
3. Tag. Don't need to create a release, it's automatically created by CI.

## Roadmap

Check out [CHANGELOG](CHANGELOG.md).
