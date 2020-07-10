📮 INBOX
=============

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/inbox)](https://goreportcard.com/report/github.com/emad-elsaid/inbox)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/inbox?status.svg)](https://godoc.org/github.com/emad-elsaid/inbox)
[![codecov](https://codecov.io/gh/emad-elsaid/inbox/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/inbox)
[![Join the chat at https://gitter.im/inbox-server/community](https://badges.gitter.im/inbox-server/community.svg)](https://gitter.im/inbox-server/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Inbox makes it easy to setup a WebRTC HTTPS signaling server

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [📮 INBOX](#📮-inbox)
    - [Install](#install)
        - [Download latest binary](#download-latest-binary)
        - [Compile from source](#compile-from-source)
        - [Docker image](#docker-image)
    - [Usage](#usage)
    - [API Documentation](#api-documentation)
    - [Purpose](#purpose)
    - [How is it working?](#how-is-it-working)
    - [The General Concept](#the-general-concept)
    - [The implementation](#the-implementation)
    - [How to run the example](#how-to-run-the-example)
    - [How to use it](#how-to-use-it)
    - [Benchmarks](#benchmarks)
    - [Contribute](#contribute)
    - [License](#license)

<!-- markdown-toc end -->

## Install

### Download latest binary

You can download [the latest version from releases](https://github.com/emad-elsaid/inbox/releases/latest) for your system/architecture

### Compile from source

- Have the [Go toolchain](https://golang.org/dl/) installed
- Clone the repository and compile and install the binary to $GOBIN
  ```
  git clone git@github.com:emad-elsaid/inbox.git
  cd inbox
  go install cmd/inbox.go
  ```

### Docker image

- If you want to run it in http mode
  ```
  docker run --rm -it -p 3000:3000 emadelsaid/inbox ./inbox --https=false
  ```
- You can use generate a self signed SSL certificate, or if you already have a
  certificate if you want to have HTTPS enabled
  ```
  docker run --rm -it -v /path/to/cert/directory:/cert -p 3000:3000 emadelsaid/inbox ./inbox --server-cert=/cert/server.crt --server-key=/cert/server.key
  ```

## Usage

```
  -bind string
        a bind for the http server (default "0.0.0.0:3000")
  -cleanup-interval int
        Interval in seconds between server cleaning up inboxes (default 1)
  -cors
        Allow CORS
  -https
        Run server in HTTPS mode or HTTP (default true)
  -inbox-capacity int
        Maximum number of messages each inbox can hold (default 100)
  -inbox-timeout int
        Number of seconds for the inbox to be inactive before deleting (default 60)
  -max-body-size int
        Maximum request body size in bytes (default 1048576)
  -max-header-size int
        Maximum request body size in bytes (default 1048576)
  -public string
        Directory path of static files to serve (default "public")
  -server-cert string
        HTTPS server certificate file (default "server.crt")
  -server-key string
        HTTPS server private key file (default "server.key")
```

## API Documentation

- Swagger documentation is under [/docs/swagger.yml](/docs/swagger.yml)
- You can show it live here https://petstore3.swagger.io/, the use the following
  URL in the top input field
  ```
  https://raw.githubusercontent.com/emad-elsaid/inbox/master/docs/swagger.yml
  ```

## Purpose

- When building a WebRTC based project you need a way to signal peers.
- One of the ways to signal peers is to use a central HTTP server
- Alice uses **Inbox** to pass WebRTC offer to Bob
- Bob gets the offer uses **Inbox** to send a WebRTC answer to Alice
- Alice and Bob use **Inbox** to exchange ICE Candidates information
- When Alice and Bob have enough ICE candidates they disconnect from **Inbox** and connect to each other directly

## How is it working?

- The server works in HTTPS mode by default unless `-https=false` passed to it.
- If you wish to generate self signed SSL certificate `server.key` and `server.crt`:
```
openssl genrsa -des3 -passout pass:secretpassword -out server.pass.key 2048
openssl rsa -passin pass:secretpassword -in server.pass.key -out server.key
openssl req -new -key server.key -out server.csr
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
```
- it uses Go to run a HTTPS server on port 3000 that serves `./public` directory
- The local server has 1 other route `/inbox` for the sender and receiver to signal each
  other the webRTC offer and answer.

## The General Concept

Inbox acts as a temporary mailbox between peers, the server creates the inbox
upon the first interaction with the user and deletes it after a duration of
inactivity which is 1 minute by default [Read more](/docs/concept.md)

## The implementation

- This is a HTTPS server written in Go
- No third party dependencies
- Stores all data in memory in one big memory structure
- Every second the server removes inboxes exceeded timeouts
- Serves `/public` in current working directory as static files
- CORS is disabled by default

## How to run the example

- Install [Go](https://golang.org/)
- Clone this repository `git clone git@github.com:emad-elsaid/inbox.git`
- Run the server `go run ./cmd/inbox.go`
- Open `https://your-ip-address:3000/send.html` on the camera machine
- Open `https://your-ip-address:3000/receive.html` on the receiver machine
- Choose the camera from the list on the sender and press `start` button
- The receiver should display the camera shortly after

## Benchmarks

Inbox inherits the high speed of Go and as it uses the RAM as storage its mostly
a CPU bound process, the project comes with go benchmarks you can test it on
your hardware yourself, You can checkout my CPU specs and benchmark numbers on
my machine here [Read more](/docs/benchmarks.md)

## Contribute

Expected contribution flow will be as follows:

* Read and understand the code
* Make some changes related to your idea
* Open a PR to validate you're in the right direction, describe what you're
  trying to do
* Continue working on your changes until it's fully implemented
* I'll merge it and release a new version

## License

MIT License (c) Emad Elsaid
