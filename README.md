# Simple gRPC user service and its CLI client

[![Build Status](https://cloud.drone.io/api/badges/maelvls/users-grpc/status.svg)](https://cloud.drone.io/maelvls/users-grpc)
[![Docker layers](https://images.microbadger.com/badges/image/maelvls/users-grpc.svg)](https://microbadger.com/images/maelvls/users-grpc)
[![Coverage Status](https://coveralls.io/repos/github/maelvls/users-grpc/badge.svg?branch=master)](https://coveralls.io/github/maelvls/users-grpc?branch=master)
[![codecov](https://codecov.io/gh/maelvls/users-grpc/branch/master/graph/badge.svg)](https://codecov.io/gh/maelvls/users-grpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/maelvls/users-grpc)](https://goreportcard.com/report/github.com/maelvls/users-grpc)
[![GolangCI](https://golangci.com/badges/github.com/maelvls/users-grpc.svg)](https://golangci.com/r/github.com/maelvls/users-grpc)
[![Godoc](https://godoc.org/github.com/maelvls/users-grpc?status.svg)](http://godoc.org/github.com/maelvls/users-grpc)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)

> So many shiny badges, I guess it doesn't mean anything anymore! I have
> been testing many services in order to select the good ones. Badges is a
> way of keeping track of them all 😅
> Conventionnal Commits and GolangCI are kind of decorative-only.

[![asciicast](https://asciinema.org/a/251067.svg)](https://asciinema.org/a/251067)

- [Simple gRPC user service and its CLI client](#simple-grpc-user-service-and-its-cli-client)
  - [Stack](#stack)
  - [Use](#use)
  - [Install](#install)
    - [Docker images](#docker-images)
    - [Binaries (Github Releases)](#binaries-github-releases)
    - [Using go-get](#using-go-get)
    - [Kubernetes & Helm](#kubernetes--helm)
  - [Develop and hack it](#develop-and-hack-it)
    - [Develop using Docker](#develop-using-docker)
  - [Technical notes](#technical-notes)
    - [Vendor or not vendor and go 1.11 modules](#vendor-or-not-vendor-and-go-111-modules)
    - [Testing](#testing)
    - [`users-cli version`](#users-cli-version)
    - [Protobuf generation](#protobuf-generation)
    - [Logs, debug and verbosity](#logs-debug-and-verbosity)
  - [Examples that I read for inspiration](#examples-that-i-read-for-inspiration)
  - [Cloud-nativeness of this project](#cloud-nativeness-of-this-project)
  - [12 factor app checklist](#12-factor-app-checklist)
  - [Kubernetes and Helm](#kubernetes-and-helm)
  - [Memo on Kubernetes](#memo-on-kubernetes)
  - [Future work](#future-work)
    - [Distributed tracing, metrics and logs](#distributed-tracing-metrics-and-logs)
    - [Service discovery and service mesh](#service-discovery-and-service-mesh)
      - [Service discovery via environement variables](#service-discovery-via-environement-variables)
      - [Service discovery via CoreDNS, Consul, Envoy or Linkerd2](#service-discovery-via-coredns-consul-envoy-or-linkerd2)
    - [Event store & event sourcing](#event-store--event-sourcing)

## Stack

- **CI/CD**: Drone.io (tests, coverage, build docker image, upload
  `users-cli` CLI binaries to Github Releases using [`goreleaser`][goreleaser])
- **Coverage**: Coveralls, Codecov
- **Code Quality**: Go Report Card, GolangCI (CI) and Pre-commit-go (local
  git hook) with:
  - **Static analysis**: gocritic, gosec, golint, goimports, deadcode,
    errcheck, gosimple, govet, ineffassign, staticcheck, structcheck,
    typecheck, unused, varcheck
  - **Formatting**: gofmt on the CI and locally with 'format on save' and
    Pre-commit-hook
- **OCI orchestration**: Kubernetes, OCI runtime = Docker, Minikube for
  testing, GKE for more testing (see related [helm-gke-terraform][])
- **Config management**: Helm
- **Dependency analysis** (the DevSecOps trend): [dependabot][] (updates go
  modules dependencies daily)
- **Local dev**: Vim, VSCode and Goland, [`gotests`][gotests],
  `golangci-lint`, `pcg` (pre-commit-go), `protoc`, `prototool`, `grpcurl`,
  `is-http2`:

  ```sh
  brew install golangci/tap/golangci-lint protobuf prototool grpcurl
  npm install -g is-http2-cli
  go install github.com/maruel/pre-commit-go/cmd/...
  ```

I created this microservice from scratch. If I was to create a new
microservice like this, I would probably use Lile for generating it (if it
needs Postres + opentracing + prom metrics + service discovery). For
example, [lile-example][].

[dependabot]: https://dependabot.com/
[gotests]: https://github.com/cweill/gotests
[lile]: https://github.com/lileio/lile
[lile-example]: https://github.com/arbarlow/account_service

## Use

Refer to [Install](#install) below for getting `users-cli` and
`users-server`.

First, let `users-server` run somewhere:

```sh
users-server
```

Then, we can query it using the CLI client. The possible actions are

- create a user
- fetch a user by his email ('get')
- list all users (the server loads some sample users on startup)
- search users by a string that matches their names
- search users by a age range

Examples:

```sh
$ users-cli create --email=mael.valais@gmail.com --firstname="Maël" --lastname="Valais" --postaladdress="Toulouse"

$ users-cli get mael.valais@gmail.com
Maël Valais <mael.valais@gmail.com> (0 years old, address: Toulouse)

$ users-cli list
Acevedo Quinn <acevedo.quinn@email.us> (22 years old, address: 403 Lawn Court, Walland, Federated States Of Micronesia, 8260)
Alford Cole <alford.cole@email.net> (33 years old, address: 763 Halleck Street, Elbert, Nevada, 3291)
Angeline Stokes <angeline.stokes@email.biz> (48 years old, address: 526 Java Street, Hailesboro, Pennsylvania, 1648)
Beasley Byrd <beasley.byrd@email.io> (56 years old, address: 213 McKibbin Street, Veguita, New Jersey, 3943)
Benjamin Frazier <benjamin.frazier@email.net> (31 years old, address: 289 Cyrus Avenue, Templeton, Maine, 5964)
Billie Norton <billie.norton@email.io> (28 years old, address: 699 Rapelye Street, Dupuyer, Ohio, 4175)
...
Stone Briggs <stone.briggs@email.info> (31 years old, address: 531 Atkins Avenue, Neahkahnie, Tennessee, 3981)
Valencia Dorsey <valencia.dorsey@email.info> (51 years old, address: 941 Merit Court, Grill, Mississippi, 4961)
Walter Prince <walter.prince@email.co.uk> (26 years old, address: 204 Ralph Avenue, Gibbsville, Michigan, 6698)
Wilkerson Mosley <wilkerson.mosley@email.biz> (48 years old, address: 734 Kosciusko Street, Marbury, Connecticut, 3037)

$ users-cli search --name=alenc
Jenifer Valencia <jenifer.valencia@email.us> (52 years old, address: 948 Jefferson Street, Guthrie, Louisiana, 2483)
Valencia Dorsey <valencia.dorsey@email.info> (51 years old, address: 941 Merit Court, Grill, Mississippi, 4961)

$ users-cli search --agefrom=30 --ageto=42
Benjamin Frazier <benjamin.frazier@email.net> (31 years old, address: 289 Cyrus Avenue, Templeton, Maine, 5964)
Stone Briggs <stone.briggs@email.info> (31 years old, address: 531 Atkins Avenue, Neahkahnie, Tennessee, 3981)
Alford Cole <alford.cole@email.net> (33 years old, address: 763 Halleck Street, Elbert, Nevada, 3291)
Brock Stanley <brock.stanley@email.me> (35 years old, address: 748 Aster Court, Elwood, Guam, 7446)
Ina Perkins <ina.perkins@email.me> (35 years old, address: 899 Miami Court, Temperanceville, Virginia, 2821)
Hardin Patton <hardin.patton@email.com> (42 years old, address: 241 Russell Street, Robinson, Oregon, 9576)
```

Here is what the help looks like:

```sh
$ users-cli help

For setting the address of the form HOST:PORT, you can
- use the flag --address=:8000
- or use the env var ADDRESS
- or you can set 'address: localhost:8000' in $HOME/.users-cli.yml

Usage:
  users-cli [command]

Available Commands:
  create      creates a new user
  get         prints an user by its email (must be exact, not partial)
  help        Help about any command
  list        lists all users
  search      searches users from the remote users-server
  version     Print the version and git commit to stdout

Flags:
      --address string   'host:port' to bind to (default ":8000")
      --config string    config file (default is $HOME/.users-cli.yaml)
  -h, --help             help for users-cli
  -v, --verbose          verbose output

Use "users-cli [command] --help" for more information about a command.
```

## Install

### Docker images

Docker images are created on each tag. The 'latest' tag represents the
latest commit on master. I use multi-stages dockerfile so that the
resulting image is less that 20MB (using Alpine/musl-libc). `latest` tag
should only be used for dev purposes as it points to the image of the
latest commit. I use [moving-tags][] `1`, `1.0` and fixed tag `1.0.0` (for
example). To run the server on port 8123 locally:

```sh
$ docker run -e LOG_FORMAT=text -e PORT=8123 -p 80:8123/tcp --rm -it maelvls/users-grpc:1
INFO[0000] serving on port 8123 (version 1.0.0)
```

[moving-tags]: http://plugins.drone.io/drone-plugins/drone-docker/#autotag

To run `users-cli`:

```sh
docker run --rm -it maelvls/users-grpc:1 users-cli --address=192.168.99.1:80 list
```

> This 172.17.0.1 address is required because communicating between
> containers through the host requires to use the IP of the docker0
> interface instead of the loopback.

### Binaries (Github Releases)

Binaries `users-cli` and `users-server` are available on the [Github
Releases page][github-releases].

Releasing binaries was not necessary (except maybe for the CLI client) but
I love the idea of Go (so easy to cross-compile + one single
statically-linked binary) so I wanted to try it. Goreleaser is a fantastic
tool for that purpose! That's where Go shines: tooling. It is exceptional
(except for [gopls][]n the Go Language Server) but it's getting better and
better). Most importantly, tooling is fast at execution and also at
compilation (contrary to Rust where compilation takes much more time --
LLVM + way richer and complex language -- see my comparison
[rust-vs-go][]).

[github-releases]: https://github.com/maelvls/users-grpc/releases
[gopls]: https://github.com/golang/go/wiki/gopls
[rust-vs-go]: https://github.com/maelvls/rust-chat

### Using go-get

```sh
go get github.com/maelvls/users-grpc/...
```

### Kubernetes & Helm

```sh
helm install ./ci/helm/users-grpc --name users-grpc --namespace users-grpc --set image.tag=latest
helm upgrade users-grpc ./ci/helm/users-grpc
```

## Develop and hack it

Here is the minimal set of things you need to get started for hacking this
project:

```sh
git clone https://github.com/maelvls/users-grpc
cd users-grpc/

brew install protobuf # only if .proto files are changed
go generate ./...     # only if .proto files are changed

go run users-server/main.go &
go run users-cli/main.go
```

### Develop using Docker

```sh
docker build . -f ci/Dockerfile --tag maelvls/users-grpc
```

In order to debug docker builds, you can stop the build process before the
bare-alpine stage by doing:

```sh
docker build . -f ci/Dockerfile --tag maelvls/users-grpc --target=builder
```

You can test the service is running correctly by using
[`grpc-health-probe`][grpc-health-probe] (note that I also ship
`grpc-health-probe` in the docker image so that liveness and readiness
checks are easy to do from kubenertes):

```sh
$ PORT=8000 go run users-server/main.go &
$ go get github.com/grpc-ecosystem/grpc-health-probe
$ grpc-health-probe -addr=:8000

status: SERVING
```

From the docker container itself:

```sh
$ docker run --rm -d --name=users-grpc maelvls/users-grpc:1
$ docker exec -i users-grpc grpc-health-probe -addr=:8000

status: SERVING

$ docker kill users-grpc
```

For building the CLI, I used the cobra cli generator:

```sh
go get github.com/spf13/cobra/cobra
```

Using Uber's [prototool][], we can debug the gRPC server (a bit like when
we use `httpie` or `curl` for HTTP REST APIs). I couple it with [`jo`][jo]
which eases the process of dealing with JSON on the command line:

```sh
$ prototool grpc --address :8000 --method user.UserService/GetByEmail --data "$(jo email='valencia.dorsey@email.info')" | jq

{
  "status": {
    "code": "SUCCESS"
  },
  "user": {
    "id": "5cfdf218f7efd273906c5b9e",
    "age": 51,
    "name": {
      "first": "Valencia",
      "last": "Dorsey"
    },
    "email": "valencia.dorsey@email.info",
    "phone": "+1 (906) 568-2594",
    "address": "941 Merit Court, Grill, Mississippi, 4961"
  }
}
```

[grpc-health-probe]: https://github.com/grpc-ecosystem/grpc-health-probe
[prototool]: https://github.com/uber/prototool
[jo]: https://github.com/jpmens/jo

## Technical notes

### Vendor or not vendor and go 1.11 modules

I use `GO111MODULES=on`. This is definitely debatable: in actual projects that I
may have to maintain on the long run, a more 'stable' option such as 'dep'
should be used (as of June 2019 at least). The reason is that I feel many
tools and libraries do not (yet) work with vgo (go 1.11 modules).

In the first iterations of this project, I was vendoring (using `go mod vendor`) and checked the vendor/ folder in with the code. Then, I realized
things have evolved and it is not necessary anymore (as of june 2019; see
[should-i-vendor][] as things may evolve).

That said, I often use `go mod vendor` which comes very handy (I can browse the
dependencies sources easily, everything is at hand).

[should-i-vendor]: https://www.reddit.com/r/golang/comments/9ai79z/correct_usage_of_go_modules_vendor_still_connects/

### Testing

I use [gotests][] for easing the TDD workflow. Whenever I add a new
function, I just have to run:

```sh
gotests -all -w users-server/service/*
```

so that these functions get generated in the corresponding `test_*.go`
file. Also, I use [go-testdeep][] in order to display a nice colorful diff
between 'got' and 'expected' for a friendlier testing experience.

[gotests]: https://github.com/cweill/gotests
[go-testdeep]: github.com/maxatome/go-testdeep

I mostly focused on TDD on `users-server`. With time, I realized that I had
many manual tests before each release. Here is a list of this that should
be added to `.drone.yml`:

1. test the docker image (at least test that the `users-server` is
   launching using [`grpc-health-probe`][grpc-health-probe])
2. test the Helm chart
3. test the CLI `users-cli` (I did not write any test for it yet)

### `users-cli version`

At build time, I use `-ldflags` for setting global variables
(`main.version` (), `main.date` (RFC3339) and `main.commit`). At first, I
was using [govvv][] to ease the process. I then realized govvv didn't help
as much as I thought; instead, if I want to have a build containing this
information, I use `-ldflags` manually (in Dockerfile for example). For
binaries puloaded to Github Releases, [`goreleaser`][goreleaser] handles
that for me. For example, a manual build looks like:

```hs
go build -ldflags "-X main.version='$(git describe --tags --always | sed 's/^v//')' -X main.commit='$(git rev-parse --short HEAD)' -X main.date='$(date --rfc-3339=date)'" ./...
```

> Note: for some reason, `-X main.date='$DATE'` cannot accept spaces in
> `$DATE` even though I use quoting. I'll have to investigate further.

[govvv]: https://github.com/ahmetb/govvv
[goreleaser]: https://github.com/goreleaser/goreleaser

### Protobuf generation

Ideally, the `.proto` and the generated `.pb.go` should be separated from
my service, e.g. `github.com/maelvls/schema` with semver versionning and
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

I did not yet implement a way for my server or my client to make the log
level higher or to set json as the logging format. For now, I use
[logrus][]. A step further (that I did not implement yet) is to log all
gRPC handlers activity (through gRPC interceptors). One way of doing that
is proposed in [go-grpc-middleware][].

## Examples that I read for inspiration

- [go-micro-services][] (lacks tests but excellent geographic-related
  business case)
- [route_guide][] (example from the official grpc-go)
- [go-scaffold][] (mainly for the BDD unit + using Ginkgo)
- [todogo][] (just for the general layout)
- [Medium: _Simple API backed by PostgresQL, Golang and
  gRPC_][medium-grpc-pg] for grpc middleware (opentracing interceptor,
  prometheus metrics, gRPC-specific logging with logrus, tags
  retry/failover, circuit-breaking -- alghouth these last two might be
  better handled by a service proxy such as [linkerd2][])
- the Go standard library was also extremely useful for learning how to
  write idiomatic code. The `net` one is a gold mine (on top of that I love
  all the networking bits).
- I learned how to publish helm charts on Github Pages there:
  [helm-gh-pages-example][]. Didn't have time to finish that part though.

[medium-grpc-pg]: https://medium.com/@vptech/complexity-is-the-bane-of-every-software-engineer-e2878d0ad45a
[go-micro-services]: https://github.com/harlow/go-micro-services
[route_guide]: https://github.com/grpc/grpc-go/tree/master/examples/route_guide
[go-scaffold]: https://github.com/orbs-network/go-scaffold
[todogo]: https://github.com/kgantsov/todogo
[helm-gh-pages-example]: https://github.com/int128/helm-github-pages
[linkerd2]: https://github.com/linkerd/linkerd2

## Cloud-nativeness of this project

Cloud-native is taking advantage when your workload is on the cloud
([cncf-definition][], [gitlab-native-talk][]):

- use containers,
- dynamically orchestrated,
- use microservices.

In this project, we use OCI containers, use Kubernetes dynamic
orchestration and use microservices.

[gitlab-native-talk]: https://youtu.be/jc5cY3LoOOI?t=204
[cncf-definition]: https://github.com/cncf/toc/blob/master/DEFINITION.md

## 12 factor app checklist

<details>
<summary>12 factor cheat sheet </summary>

> Source: [12factor-list][]

| ✓   | Factors                                                                    | Status                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| --- | -------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Codebase                                                                   | One codebase tracked in revision control, many deploys. One Code base one repo is handling all the environment. ex: production, staging, integration, local                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| 2   | Dependencies: Explicitly declare and isolate dependencies                  | All the dependencies are declaraed outside the CodeBase. Pip installable library is used and virtualenv for isolation of dependencies.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| 3   | Config: store config in the environment                                    | strict separation of config from code. Config varies substantially across deploys, code does not. Store config in environment variables. Will store env variables in env file not tracked and used by the code, used to declare environment variables only.                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| 4   | Backing services: treat backing services as attached resources             | A backing service is any service the app consumes over the network as part of its normal operation. Examples include MySQL, RabbitMQ both are attached resources, accessed via a URL or other locator/credentials stored in the config(ENV_VAR). Attached resources must be change without any code changes. **Remark**: depends on 3.                                                                                                                                                                                                                                                                                                                                                                                        |
| 5   | Build, release, run: strictly separate build and run stages                | We have to maintain three stages to release a project: Build stage(Convert a project to build module using the executable commit from version system it fetches vendors dependencies and compiles binaries and assets.) Release stage (takes the build produced by the build stage and combines it with the deploy’s current config). Run stage(Runs the app in the execution environment, by launching some set of the app’s processes against a selected release using the gunicorn, worker or the supervisor). Have to use some deployment tools so that every release should have a release ID and having the capability to rollback to a particular release ID. Docker-based containerized deploy strategy would be used |
| 6   | Processes: execute the app as one or more stateless processes              | App should be stateless and any data that needs to persist should be managed through a stateful backing service like mysql and rabbitMq.The Memory space or filesystem of the process can be used as a brief, single-transaction cache, So that it can be run through multiple processes. And gunicorn maintaining the one or more stateless processes                                                                                                                                                                                                                                                                                                                                                                        |
| 7   | Port binding: export services via port binding                             | The web app exports HTTP as a service by binding to a port, and listening to requests coming in on that port. Example: the developer visits a service URL like <http://localhost:5000/> to access the service exported by their app. Running the flask app through gunicorn and bind it to IP and PORT which you want to use.                                                                                                                                                                                                                                                                                                                                                                                                 |
| 8   | Concurrency: scale out via the process model                               | Architect app to handle diverse workloads by assigning each type of work to a process type. For example, HTTP requests may be handled by a web process, and long-running background tasks handled by a worker process. Application must also be able to span multiple processes running on multiple physical machines.                                                                                                                                                                                                                                                                                                                                                                                                        |
| 9   | Disposability: maximize robustness with fast startup and graceful shutdown | Processes should strive to minimize startup time. Ideally, a process takes a few seconds from the time the launch command is executed until the process is up and ready to receive requests or jobs. Short startup time provides more agility for the release process and scaling up; and it aids robustness, because the process manager can more easily move processes to new physical machines when warranted. and a graceful shutdown. Flask server should be shutdown with supervisor stop as it makes the process to shutdown gracefully                                                                                                                                                                                |
| 10  | Dev/prod parity: keep dev, staging, and prod as similar as possible        | Make the time gap small: a developer may write code and have it deployed hours or even just minutes later. Make the personnel gap small: developers who wrote code are closely involved in deploying it and watching its behavior in production.Make the tools gap small: keep development and production as similar as possible. **Remark:** redeploy every hour, code authors must be the deployers, dev and prod env must be as same as possible                                                                                                                                                                                                                                                                           |
| 11  | Logs: treat logs as event streams                                          | App should not attempt to write to or manage logfiles. Instead, each running process writes its event stream, unbuffered, to stdout.In staging or production deploys, each process’ stream will be captured by the execution environment, collated together with all other streams from the app, and routed to one or more final destinations for viewing and long-term archival. Should Follow ELK.                                                                                                                                                                                                                                                                                                                          |
| 12  | Admin processes: run admin/management tasks as one-off processes           | Any admin or management tasks for a 12-factor app should be run as one-off processes within a deploy’s execution environment. This process runs against a release using the same codebase and configs as any process in that release and uses the same dependency isolation techniques as the long-running processes.                                                                                                                                                                                                                                                                                                                                                                                                         |

</details>

[12factor][] is a manifest originally proposed by Heroku and largely
adopted among the cloud-native ommunity. It presents the 12 main ideas that
should be thought of when building an application that is meant to be run
on a cloud provider (e.g., platforms like Now.sh or Heroku or any other
cloud-oriented platform such as Kubernetes).

Here is a checklist for my microservice and its CLI:

1.  **Codebase:**

    In this repo, the master represents the branch continously deployed to
    the staging environment.

    - features are developed using pull requests (as they show in one place
      the code-review process as well as the tests & coverage status).
    - once merged, the PR automatically triggers a deployment of the
      service on the staging environment.
    - master should always be in a deployable state.
    - tags are 'immutable and identifiable releases' that trigger a
      deployment to the production env (might need a bit more process here,
      like having a `prod` branch and PRs from `master` to `prod`; I found
      [gitlab-workflow][] very interesting on that subject).

    Note: (this is an opinion!) I don't think that another branch like `staging`
    or `develop` should be introduced unless a pull-request workflow is needed
    (like for the `production` branch at Gitlab) in which case it becomes
    necessary. Git-flow is nice but too bloated (not very KISS 😁).

    [gitlab-workflow]: https://about.gitlab.com/handbook/engineering/infrastructure/design/git-workflow/
    [gitlab-handbook]: https://about.gitlab.com/handbook

2.  **Dependencies: Explicitly declare and isolate dependencies:**

    In this project, I use 'go modules' (go 1.11) which uses `go.sum` for
    locking dependencies, promoting reproducible builds. One exception:
    `protoc`, the protobuf generator, is not version locked but is only
    needed when modifying .proto files (I didn't find a workaround on that
    issue yet).

3.  **Config: store config in the environment:**

    The server is confugurable using env vars: `PORT` (defaults to 8000) and
    `LOG_FORMAT=text|json` (`text` mode by default). The docker images set
    sensible defaults for these env vars.

4.  **Backing services: treat backing services as attached resources:**

    Although I don't have any DB, redis or event store or external API calls
    in this project. In case of resource changes, everything would be
    configured using env vars; any change of env var would require to
    relaunch the service though (which seems to be the correct way of
    doing).

    Alternatively, I could also use Consul instead of env vars for passing
    other services ip/port and credentials (or even better: Hashicorp
    Vault). My server would query the credentials on startup. It has nice
    advantages but also means more logic into each microservice and also
    being tied to a specific 'way', compared to generic and pervasive env
    vars.

5.  **Build, release, run: strictly separate build and run stages:**

    Build and Release (build + config) are separated

    - build = docker image
    - release = build (docker image) + config (Deployment + Service) on
      kubernetes. Release ID = the history ID given by:

      ```sh
      kubectl rollout history deployment.v1.apps/users-grpc
      ```

      And we can rollback using `kubectl rollout undo`.

      Similarly to Kubernetes, Helm provides a nice fallback mechanism: when
      `helm upgrade users-grpc` fails, it will fall back the last known state
      of this helm release.

    - run = use Kubernetes as the supervisor, run the container in a Pod
      thanks to the Deployment config.

6.  **Processes: execute the app as one or more stateless processes:**

    Using the OCI runtime, memory spaces and filesystems are isolated using
    namespaces. This service is not stateless (as I did not have the time
    for using a DB such as Postgres + go-pg because microservice = small so
    probably no need for ORM but maybe I need a migration tool).

    Regarding resource limiting, I do not have limits or requests in this
    service (minikube). In a real-case scenario, resource requests would
    allow a depending process to be guaranteed to function under correct
    circumstances.

7.  **Port binding: export services via port binding:**

    `PORT` is available for setting the port. It uses L3 TCP protocol and
    HTTP/2 + gRPC for L5 to L7.

8.  **Concurrency: scale out via the process model:**

    OCI containers help a lot with that; in a microservice architecture,
    worker processes (long-term) and web services (short-term) would
    probably not be grouped in the same microservice as they have different
    concerns.

9.  **Disposability: maximize robustness with fast startup and graceful shutdown:**

    Fast startup: using Go and no DB, startup takes less that 20ms. Using a
    DB, it would still be something like 1 second, which definitely counts
    as a 'short startup'.

    Regarding the 'graceful shutdowns', meh... I did not (yet) study how
    Kubernetes is handling container kills and such.

10. **Dev/prod parity: keep dev, staging, and prod as similar as possible:**

    Here, the OCI image is the same from the dev environement (local) to
    the prod environement.

    By the way, the idea that

    > The code author must be the one who deploys

    means that developer should be involved with deployment concerns, and
    that the best person for the deployment task is the developer himself.
    That idea works well with the PR being merged into master (i.e., a
    deploy to the staging env). But who is in charge of tagging (and thus,
    deploying into production)?

11. **Logs: treat logs as event streams:**

    This microservice prints JSON (or plain text) logs to stdout. I did not
    log the method calls, but it is possible to do it using gRPC
    interceptors (just one line of Go).

    Tracing is not addressed in the 12factor recommandations and I think
    there should exist a '12 factor kube' where metrics and distributed
    tracing is addressed.

    Service discovery and service mesh (Consul, [linkerd2][], Istio) isn't
    addressed either.

12. **Admin processes: run admin/management tasks as one-off processes:**

    Admin tasks, such as debug or database operations, can be handled using
    `docker exec` inside the running container, which means it shares the
    same release (build + config).

[12factor]: http://12factor.net
[12factor-list]: https://gist.github.com/anandtripathi5/118995139602599dab64fddcd147545a

## Kubernetes and Helm

In order to test the deployment of my service, I create a Helm chart (as
well as a static `kubernetes.yml` -- which is way less flexible) and used
minikube in order to test it. I implemented the [grpc-healthcheck][] so that Kubernetes's readyness and
liveness checks can work with this service. What I did:

1. clean logs in JSON ([logrus][]) for easy integration with Elastic/ELK
2. health probe working (readiness)
3. `helm test --cleanup users-grpc` passes
4. the service can be exposed via an Ingress controller such as Traefik or
   Nginx. For example, using the Helm + GKE + Terraform configuration at
   [helm-gke-terraform][]:

   ```yaml
   image:
     tag: 1.0.0
   ingress:
     enabled: true
     hosts: [users-grpc.kube.maelvls.dev]
     annotations:
       kubernetes.io/ingress.class: traefik
       certmanager.k8s.io/cluster-issuer: letsencrypt-prod
     tls:
       - hosts: [users-grpc.kube.maelvls.dev]
         secretName: users-grpc-example-tls
   ```

   We can then have the service from the internet through Traefik (Ingress
   Controller) with dynamic per-endpoint TLS ([cert-manager][]) and DNS
   ([external-dns][]):

   ```sh
   helm install ./helm/users-grpc --name users-grpc --namespace users-grpc --set image.tag=latest --values helm/users-grpc.yaml
   ```

[helm-gke-terraform]: https://github.com/maelvls/awx-gke-terraform
[grpc-healthcheck]: https://github.com/grpc/grpc/blob/master/doc/health-checking.md
[logrus]: https://github.com/sirupsen/logrus
[external-dns]: https://github.com/kubernetes-incubator/external-dns
[cert-manager]: https://github.com/jetstack/cert-manager

To bootstrap the kubernetes YAML configuration for this service using my
Helm chart, I used:

```sh
helm template ./ci/helm/users-grpc --name users-grpc --namespace users-grpc --set image.tag=latest > ci/deployment.yml
```

We can now apply the configuration without using Helm. Note that I changed
the ClusterIP to NodePort so that no LoadBalancer, Ingress Controller nor
`kubectl proxy` is needed to access the service.

```sh
$ kubectl apply -f ci/deployment.yml
$ kubectl get svc users-grpc
NAME        TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
users-grpc   NodePort   10.110.71.154   <none>        8000:32344/TCP   15m
```

Now, in order to access it, we must retrieve the minikube cluster IP (i.e.,
its service IP, the IP used by kubectl for sending commands).

```sh
$ minikube status
kubectl: Correctly Configured: pointing to minikube-vm at 192.168.99.105
```

We then use the address `192.168.99.105:32344`. Let's try with
[grpc-health-probe][]:

```sh
% grpc-health-probe -addr=192.168.99.105:32344
status: SERVING
```

Yey!! 🎉🎉🎉

## Memo on Kubernetes

- Controllers

  - **deployment** (for scaling) -> ReplicaSet (services + nb of replicas) -> Pod (containers)
  - **StatefulSet** (for volumes) -> Pods -> PersistentVolumeClaim
  - **Job** (one-shot op)
  - **DeamonSet** (one deamon per node) -> for node-exporter and filebeat

- ClusterIP vs NodePort vs Ingress vs LoadBalancer:

  - **ClusterIP** exposes a Service. It is the default for a Service -> no
    external access (must use an Ingress for that Service)
  - **NodePort** exposes a Service. It redirects requests for a given port to a Service ->
    easy external access but not very good way of exposing services. Good
    for healthchecks; I also use it for the ACME cert challenge with
    `certmanager`.
  - **LoadBalancer** also kind of exposes a Service (the load balancer
    one). It is like ClusterIP but for exposing to the internet.
  - **Ingress** -> not a service; it maps a Service to an Ingress
    Controller.

- **Ingress Controller** (Traefik or Nginx or Cloud-specific LBs) applies
  Ingress (maps a Service to an Ingress Controller)
- **ConfigMap** -> for storing config files like `traefik.toml`
- **Secret** -> for storing credentials and certificates

## Future work

Here is a small list of things that could be implemented now that a MVP
microservice is working.

### Distributed tracing, metrics and logs

- Prometheus: one way would be to use a grpc interceptor that would send
  metrics to prometheus
- Jaeger: very nice for debugging a cascade of gRPC calls. It requires a
  gRPC interceptor compatible with Opentracing.
- Logs: [logrus][] can log every request or only failing requests, and this can
  be easily implemented using a gRPC interceptor (again!)

These middlewares are listed and available at [go-grpc-middleware][].

[go-grpc-middleware]: https://github.com/grpc-ecosystem/go-grpc-middleware

### Service discovery and service mesh

How can other services use it from inside the cluster? As stated in the
documentation ([connect-applications-service][]),

> Kubernetes supports 2 primary modes of finding a Service - environment
> variables and DNS. The former works out of the box while the latter
> requires the CoreDNS cluster addon.

#### Service discovery via environement variables

Let's try the env var approach (using minikube):

```sh
$ kubectl get pods
NAME                               READY   STATUS    RESTARTS   AGE
users-grpc-69d46c866f-t6rx4         1/1     Running   0          23m

$ kubectl exec -it users-grpc-69d46c866f-t6rx4 env | grep -i users-grpc
HOSTNAME=users-grpc-69d46c866f-t6rx4
...
USERS_GRPC_SVC_SERVICE_PORT=8000
USERS_GRPC_SVC_SERVICE_HOST=10.110.71.154
```

Let us say we have service A that wants to use users-grpc. Service A will be
provided with these env variables. Note that because of the dependency on
users-grpc, this service would probably fail on startup until users-grpc is
up. Requires some extra logic on startup.

#### Service discovery via CoreDNS, Consul, Envoy or Linkerd2

Service discovery can also directly use the Kubernetes API from the service
itself, or using a sidekick container ([linkerd2][] or [envoy][] as a
service proxy) or with a [grpc-consul-resolver][]. [linkerd2][] and
[envoy][] also add the possibility of circuit breaking and service-level
(L7) load-balancing.

[connect-applications-service]: https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/
[envoy]: https://www.envoyproxy.io/
[grpc-consul-resolver]: https://github.com/mbobakov/grpc-consul-resolver

### Event store & event sourcing

In case this README didn't have enough buzzwords, let's add an event store
to the stack. Instead of a traditional DB (which only stores the current
state), an event store keeps tracks of all RPC messages that have been sent
to him; the current state is then re-computed every time the service needs
to access it. We call this model 'event sourcing'. It requires fast reading
speed, which NoSQL DBs such as Mongodb excel at.

Event stores and event sourcing often also implies an event bus such as
Kafka (or RabbitMQ) so that other services can register to topics; for
example:

- service `User` receives an event (= message) `add user Claudia Greene`
- service `User` processes the message and adds the user
- service `User` sends a new message to the topic 'user-added'
- service `Metrics` is subscribing to 'user-added'; it receives message
  that says that a user has been added and adds +1 to its count of new
  users for the month.
- service `Email` is also subscribing to the topic 'user-added'. Upon
  reception, it sends an email to Claudia Greene for welcoming her.
