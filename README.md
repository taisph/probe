# Probe

[![Go Report Card](https://goreportcard.com/badge/github.com/taisph/probe)](https://goreportcard.com/report/github.com/taisph/probe)

Probe is a simple tool to wait for access to network ports and sockets.

This tool is intended to be used in a Docker container right before launching the actual application. It helps ensure
there is access to the network hosts or unix sockets the application requires, before the application is launched. This
is especially usefull when using Kubernetes with a service mesh where network access may be blocked for a short period
while proxy sidecars start up.

# Usage

Probe currently supports TCP and Unix socket addresses.

Wait for localhost port 8080, google.com port 443 and the MySQL socket:

```shell
probe wait localhost:8080 tcp:google.com:443 unix:/var/run/mysqld.sock
```

To see more options, like logging and timeouts:

```shell
probe wait --help
```

# Installation

## Precompiled

There are procompiled binaries under [releases](https://github.com/taisph/probe/releases).

## Source

Requires [go](https://golang.org/).

To retrieve, build and install the command, run:

```bash
go get -u github.com/taisph/probe/cmd/probe
```
