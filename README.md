# Simple gRPC quote service and its nice CLI

[![Build
Status](https://cloud.drone.io/api/badges/maelvalais/quote/status.svg)](https://cloud.drone.io/maelvalais/quote)
[![Coverage
Status](https://coveralls.io/repos/github/maelvalais/quote/badge.svg?branch=master)](https://coveralls.io/github/maelvalais/quote?branch=master)
[![codecov](https://codecov.io/gh/maelvalais/quote/branch/master/graph/badge.svg)](https://codecov.io/gh/maelvalais/quote)
[![GolangCI](https://golangci.com/badges/github.com/maelvalais/quote.svg)](https://golangci.com/r/github.com/maelvalais/quote)
[![Godoc](https://godoc.org/github.com/maelvalais/quote?status.svg)](http://godoc.org/github.com/maelvalais/quote)
[![Go Report Card](https://goreportcard.com/badge/github.com/maelvalais/quote)](https://goreportcard.com/report/github.com/maelvalais/quote)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)

> Sorry for the many badges, I have been testing many services in order to
> select the good ones. 😅 Conventionnal Commits and GolangCI are kind of
> decorative-only.

## Install

### Dev

```sh
brew install protobuf # only if .proto files are changed
go generate ./...     # only if .proto files are changed
go build
./quote
```

### Docker

```sh
docker build . -f ci/Dockerfile
```

For building the CLI, I used the cobra cli generator:

```sh
go get github.com/spf13/cobra/cobra
```

## Side notes

### Vendor or not vendor and go 1.11 modules

I use `GO111MODULES=on`. This is definitely debatable: in actual projects that I
may have to maintain on the long run, a more 'stable' option such as 'dep'
should be used (as of June 2019 at least). The reason is that I feel many
tools and libraries do not (yet) work with vgo (go 1.11 modules).

In the first iterations of this project, I was vendoring (using `go mod vendor`)
and checked the vendor/ folder in with the code. Then, I realized things have
evolved and it is not necessary anymore (as of june 2019; see [should-i-vendor]
as things may evolve).

That said, I often use `go mod vendor` which comes very handy (I can browse the
dependencies sources easily, everything is at hand).

[should-i-vendor]: https://www.reddit.com/r/golang/comments/9ai79z/correct_usage_of_go_modules_vendor_still_connects/

### `quote version`

I decided to use <https://github.com/ahmetb/govvv> in order to ease the
process of using `-ldflags -Xmain.Version=$(git describe)` and so on. I
could have done it without it 🙄

### Proto generation

Ideally, the `.proto` and the generated `.pb.go` should be separated from
my service, e.g. `github.com/maelvalais/schema` with semver versionning and
auto-generated `.pb.go` by the CI (see this [SO
discussion](proto-monorepo)). Or maybe the `.pb.go` should be owned by
their respective services... Depending on the use of GO111MODULES or `dep`.

For `*.pb.go` generation, I use the annotation `//go:generate protoc`. In order
to re-generate the pb files from the proto files, don't forget to do:

```sh
go generate ./...
```

[proto-monorepo]: https://stackoverflow.com/questions/55250716/organization-of-protobuf-files-in-a-microservice-architecture

### Logs, debug and verbosity

I did not yet implement a way for my server or my client to make the log level
higher or to set json as the logging format. For now, I use logrus; fortunately,
logrus allows to configure these things.

<!--
I did not implement a way of logrotating the logs ([traefik's log rotation][traefik-logrotate]
is an excellent source of inspiration in that regard)

[traefik-logrotate]: https://docs.traefik.io/configuration/logs/#log-rotation -->

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

## Cloud native

Cloud-native is taking advantage when your workload is on the cloud ([cncf-definition], [gitlab-native-talk]):

- use containers,
- dynamically orchestrated
- use micro-services

[gitlab-native-talk]: https://youtu.be/jc5cY3LoOOI?t=204
[cncf-definition]: https://github.com/cncf/toc/blob/master/DEFINITION.md

## 12factor

[12factor] is a manifest originally proposed by Heroku and largely adopted among
the community. It presents the 12 main ideas that should be thought of when
building an application that is meant to be run on a cloud provider (e.g.,
platforms like Now.sh or Heroku or any other cloud-oriented platform such as
Kubernetes). Here is a checklist for my microservice and its CLI (source:
[12factor-list]):

| ✓   | Factors                                                                           | Status                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | Remarks                                                                                                   |
| --- | --------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| 1   | CodeBase                                                                          | One codebase tracked in revision control, many deploys. One Code base one repo is handling all the environment ex: production, staging, integration, local                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |                                                                                                           |
| 2   | Dependencies(Explicitly declare and isolate dependencies)                         | All the dependencies are declaraed outside the CodeBase. Pip installable library is used and virtualenv for isolation of dependencies.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |                                                                                                           |
| 3   | Config(Store config in the environment)                                           | strict separation of config from code. Config varies substantially across deploys, code does not. Store config in environment variables. Will store env variables in env file not tracked and used by the code, used to declare environment variables only.                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |                                                                                                           |
| 4   | Backing services(Treat backing services as attached resources)                    | A backing service is any service the app consumes over the network as part of its normal operation. Examples include MySQL, RabbitMQ both are attached resources, accessed via a URL or other locator/credentials stored in the config(ENV_VAR). Attached resources must be change without any code changes.(DEPENDENCY ON 3)                                                                                                                                                                                                                                                                                                                                                                                                 | (Dependency on 3)                                                                                         |
| 5   | Build, release, run(Strictly separate build and run stages)                       | We have to maintain three stages to release a project: Build stage(Convert a project to build module using the executable commit from version system it fetches vendors dependencies and compiles binaries and assets.) Release stage (takes the build produced by the build stage and combines it with the deploy’s current config). Run stage(Runs the app in the execution environment, by launching some set of the app’s processes against a selected release using the gunicorn, worker or the supervisor). Have to use some deployement tools so that every release should have a release ID and having the capability to rollback to a particular release ID.Docker-based containerized deploy strategy would be used |                                                                                                           |
| 6   | Processes(Execute the app as one or more stateless processes)                     | App should be stateless and any data that needs to persist should be managed through a stateful backing service like mysql and rabbitMq.The Memory space or filesystem of the process can be used as a brief, single-transaction cache, So that it can be run through multiple processes. And gunicorn maintaining the one or more stateless processes                                                                                                                                                                                                                                                                                                                                                                        |                                                                                                           |
| 7   | Port binding(Export services via port binding)                                    | The web app exports HTTP as a service by binding to a port, and listening to requests coming in on that port. Example: the developer visits a service URL like http://localhost:5000/ to access the service exported by their app. Running the flask app through gunicorn and bind it to IP and PORT which you want to use.                                                                                                                                                                                                                                                                                                                                                                                                   |                                                                                                           |
| 8   | Concurrency(Scale out via the process model)                                      | Architect app to handle diverse workloads by assigning each type of work to a process type. For example, HTTP requests may be handled by a web process, and long-running background tasks handled by a worker process. Application must also be able to span multiple processes running on multiple physical machines.                                                                                                                                                                                                                                                                                                                                                                                                        |                                                                                                           |
| 9   | Disposability(Maximize robustness with fast startup and graceful shutdown)        | Processes should strive to minimize startup time. Ideally, a process takes a few seconds from the time the launch command is executed until the process is up and ready to receive requests or jobs. Short startup time provides more agility for the release process and scaling up; and it aids robustness, because the process manager can more easily move processes to new physical machines when warranted. and a graceful shutdown. Flask server should be shutdown with supervisor stop as it makes the process to shutdown gracefully                                                                                                                                                                                |                                                                                                           |
| 10  | Dev/prod parity(Keep development, staging, and production as similar as possible) | Make the time gap small: a developer may write code and have it deployed hours or even just minutes later.Make the personnel gap small: developers who wrote code are closely involved in deploying it and watching its behavior in production.Make the tools gap small: keep development and production as similar as possible.                                                                                                                                                                                                                                                                                                                                                                                              | deploy time: hours, code authors and deployers: same, Dev and production environment: as same as possible |
| 11  | Logs(Treat logs as event streams)                                                 | App should not attempt to write to or manage logfiles. Instead, each running process writes its event stream, unbuffered, to stdout.In staging or production deploys, each process’ stream will be captured by the execution environment, collated together with all other streams from the app, and routed to one or more final destinations for viewing and long-term archival. Should Follow ELK.                                                                                                                                                                                                                                                                                                                          |                                                                                                           |
| 12  | Admin processes(Run admin/management tasks as one-off processes)                  | Any admin or management tasks for a 12-factor app should be run as one-off processes within a deploy’s execution environment. This process runs against a release using the same codebase and configs as any process in that release and uses the same dependency isolation techniques as the long-running processes.                                                                                                                                                                                                                                                                                                                                                                                                         |                                                                                                           |

[12factor]: http://12factor.net
[12factor-list]: https://gist.github.com/anandtripathi5/118995139602599dab64fddcd147545a

## Go popularity

When I was learning Rust, I did a short 'Go vs Rust': [rust-vs-go]. The gist is that...

[rust-vs-go]: https://github.com/maelvalais/rust-chat
