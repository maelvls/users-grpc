# Simple gRPC quote service and its nice CLI

[![Build Status](https://cloud.drone.io/api/badges/maelvalais/quote/status.svg)](https://cloud.drone.io/maelvalais/quote)

Notes:

1. I use the `// +build tools` convention which allows me to specify dev
   dependencies (like protoc). See:
   <https://github.com/golang/go/issues/25922>.
2. I use `//go:generate protoc` so I run `go generate` whenever I update my
   proto files.
3. I use `GO111MODULES=on`. This is definitely debatable; I also use `go mod vendor` which comes very handy (I can browse the dependencies
   sources easily). Other approaches: use `dep` and stay in `$GOPATH/src`.
4. For building a CLI in 10 minutes, I used go-cli:

   ```sh
   GO111MODULE=off go get -d github.com/tcnksm/gcli && (cd \$(go env GOPATH)/src/github.com/tcnksm/gcli && make install)
   go get github.com/spf13/cobra/cobra
   ```

## Install

```sh
brew install protobuf # only if .proto files are changed
go generate ./...     # only if .proto files are changed
go build
./quote
```

## Side notes

### Vendor or not vendor

At first, I use `go mod vendor` and checked the vendor/ folder in with the
code. Then, I realized things have evolved and it is not necessary anymore
(as of june 2019; see [should-i-vendor] as things may evolve).

[should-i-vendor]: https://www.reddit.com/r/golang/comments/9ai79z/correct_usage_of_go_modules_vendor_still_connects/

### `quote version`

I decided to use <https://github.com/ahmetb/govvv> in order to ease the
process of using `-ldflags -Xmain.Version=$(git describe)` and so on. I
could have done without it.

### Proto generation

Ideally, the `.proto` and the generated `.pb.go` should be separated from
my service, e.g. `github.com/maelvalais/schema` with semver versionning and
auto-generated `.pb.go` by the CI (see this [SO
discussion](proto-monorepo)). Or maybe the `.pb.go` should be owned by
their respective services... Depending on the use of GO111MODULES or `dep`.

[proto-monorepo]: https://stackoverflow.com/questions/55250716/organization-of-protobuf-files-in-a-microservice-architecture

### Logs

I did not implement a way of logrotating the logs ([traefik's log rotation][traefik-logrotate]
is an excellent source of inspiration in that regard)

[traefik-logrotate]: https://docs.traefik.io/configuration/logs/#log-rotation

### Static analysis, DevSecOps and CI

- CI and commit hook : <https://github.com/golangci/golangci-lint> which
  runs gosec, go-critic and so on.
- CI and commit hook: conventional-commit using github.com/geocine/golem
  (see [golem-post]). I also like using emojis
  (github.com/carloscuesta/gitmoji-cli) but the use of emojis is much more
  debatable 😁
- `GO111MODULES=off go get github.com/maruel/pre-commit-go/cmd/...` and
  `pcg install` (caveat: it forces the dev to have no untracked file before
  committing... good thing from my experience: it forces you to check in
  the untracked files as soon as possible, see [pcg-untracked-issue]). Pcg
  will run go build on every commit so that you don't have a failing commit
  (easier to `git bisect` 👍)

[pcg-untracked-issue]: https://github.com/maruel/pre-commit-go/issues/15
[golem-post]: https://dev.to/erinbush/being-intentional-with-commits--59a3

## Examples that I read for inspiration

- [go-micro-services] (lacks tests but excellent geographic-related
  business case)
- [route_guide] (example from the official grpc-go)
- [go-scaffold] (mainly for the BDD unit + using Ginkgo)
- [todogo] (just for the general layout)

[go-micro-services]: https://github.com/harlow/go-micro-services
[route_guide]: https://github.com/grpc/grpc-go/tree/master/examples/route_guide
[go-scaffold]: https://github.com/orbs-network/go-scaffold
[todogo]: https://github.com/kgantsov/todogo

## Tools I used

- editors: vim, vscode and goland
- misc: prototool (for testing the gRPC server)
