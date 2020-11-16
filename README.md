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

- [Stack](#stack)
- [Use](#use)
- [Install](#install)
  - [Docker images](#docker-images)
  - [Binaries (Github Releases)](#binaries-github-releases)
  - [Using go-get](#using-go-get)
  - [Kubernetes & Helm](#kubernetes--helm)
- [Develop and hack it](#develop-and-hack-it)
  - [Testing](#testing)
  - [Develop using Docker](#develop-using-docker)
- [Technical notes](#technical-notes)
  - [Vendor or not vendor and go 1.11 modules](#vendor-or-not-vendor-and-go-111-modules)
  - [Testing](#testing-1)
  - [`users-cli version`](#users-cli-version)
  - [Protobuf generation](#protobuf-generation)
  - [Logs, debug and verbosity](#logs-debug-and-verbosity)
- [Examples that I read for inspiration](#examples-that-i-read-for-inspiration)
- [Kubernetes and Helm](#kubernetes-and-helm)
- [Future work](#future-work)
  - [Distributed tracing, metrics and logs](#distributed-tracing-metrics-and-logs)

## Stack

- **CI/CD**: Drone.io (tests, coverage, build docker image, upload
  `users-cli` CLI binaries to Github Releases using [`goreleaser`][goreleaser])
- **Coverage**: Coveralls, Codecov
- **Code Quality**: Go Report Card, GolangCI (CI & local git pre-push
  hook).
- **OCI orchestration**: Kubernetes,
  [Kind](https://github.com/kubernetes-sigs/kind) for testing, Civo for
  live testing (see related [k.maelvls.dev][])
- **Config management**: Helm
- **Dependency analysis** (the DevSecOps trend): [dependabot][] (updates go
  modules dependencies daily)
- **Local dev**: Vim, VSCode and Goland, [`gotests`][gotests],
  `golangci-lint`, `pcg` (pre-commit-go), `protoc`, `prototool`, `grpcurl`:

  ```sh
  brew install golangci/tap/golangci-lint protobuf prototool grpcurl
  ```

[dependabot]: https://dependabot.com/
[gotests]: https://github.com/cweill/gotests

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

go run ./users-server &
go run ./users-cli
```

### Testing

I wrote two kinds of tests:

- Unit tests to make sure that the database logic works as expected. Tests
  are wrapped in transactions which are rolled back after the test. To run the unit tests:

  ```sh
  go test ./... -short
  ```

- End-to-end tests where both the CLI and server are built and run. These
  tests check the user-facing behaviors, e.g., that the CLI arguments work
  as expected and that the CLI returns the expected exit code. To run those:

  ```sh
  go test ./test/e2e
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
$ PORT=8000 go run ./users-server &
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
go build -ldflags "-X main.version='$(git describe --tags --always | sed 's/^v//')' -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date --rfc-3339=date)" ./...
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

## Kubernetes and Helm

In order to test the deployment of my service, I create a Helm chart (as
well as a static `kubernetes.yml` -- which is way less flexible) and used
minikube in order to test it. I implemented the [grpc-healthcheck][] so that Kubernetes's readyness and
liveness checks can work with this service. What I did:

1. clean logs in JSON ([logrus][]) for easy integration with Elastic/ELK
2. health probe working (readiness)
3. `helm test --cleanup users-grpc` passes
4. the service can be exposed via an Ingress controller such as Traefik or
   Nginx. For example, using the Helm + Civo K3s + Terraform configuration at
   [k.maelvls.dev][]:

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

[k.maelvls.dev]: https://github.com/maelvls/k.maelvls.dev
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
