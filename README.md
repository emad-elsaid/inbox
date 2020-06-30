📮 INBOX
=============

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/inbox)](https://goreportcard.com/report/github.com/emad-elsaid/inbox)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/inbox?status.svg)](https://godoc.org/github.com/emad-elsaid/inbox)
[![codecov](https://codecov.io/gh/emad-elsaid/inbox/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/inbox)
[![Join the chat at https://gitter.im/inbox-server/community](https://badges.gitter.im/inbox-server/community.svg)](https://gitter.im/inbox-server/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


## Purpose

- When building a WebRTC based project you need a way to signal peers.
- One of the ways to signal peers is to use a central HTTP server
- The initiating peer use the server to pass WebRTC offer to the dialed peer
- The dialed peer gets the offer and send an answer to the initiating peer
- Then peers exchange ICE Candidates information through the server until they have enough ICE candidates to connect to each other
- The server in this repository acts as a mailbox for peers to exchange these previous offers/answers/ice candidates

## How is it working?

- the repository includes a script to generate SSL self signed certificate
  `ssl-gen` as it's needed to run the server and use webRTC in development/locally
- it uses Go to run a HTTPS server on port 3000 that serves `public` directory
- The local server has 1 other route `/inbox` for the sender and receiver to signal each
  other the webRTC offer and answer.

## The General Concept

- The server acts as temporary mailbox for peers
- Peers use basic authentication (username, password) to get or send messages
- Whenever a peer authenticate with username and password an inbox will be
  created for them if it doesn't exist
- Username can be anything: random number, UUID, public key...etc
- If the username exists and the password is correct then the server will
  respond with the oldest message in the inbox and deletes it from it's memory,
  and will respond with header `X-From` with the peer username that sent this
  message.
- If the username exists and the password is incorrect an Unauthorized arror is returned
- Now the Inbox with this username is ready to receive messages from another peer.
- A peer can use another peer username to push a message to his inbox
- The peer inbox will expire after a period of time (1 minute by default) of not
  asking for any message
- The message has a timeout and will be deleted after this timeout (1 minute by default)
- So peers has to keep asking the server for new messages with short delays that
  doesn't exceed the timeout until they got enough information to connect to
  each other
- So for 2 peers to connect, the first peer need to choose an identifier and
  pass it to the other peer in any other medium (Chat or write it on a paper or
  pre share it)
- The first peer use it to create his inbox and wait for messages from any peer
- The second peer will create an inbox with any username and send a message to
  initiate connect to the pre shared username.

## The implementation

- This is a HTTPS server written in Go
- Stores all data in memory in one big memory structure
- Clears old data every second to remove old inboxes and messages
- It started as a backend for sharing my phone camera with my PC and the idea evolved to cover more usecases, this is why the example in `public` shares the camera.

## How to run the example

- Clone the code
- make sure you have Go installed
- run the server `go run ./cmd/inbox.go`
- open `https://your-ip-address:3000/send.html` on the camera machine
- open `https://your-ip-address:3000/receive.html` on the receiver machine
- choose the camera from the list on the sender and press `start` button
- the receiver should display the camera shortly after

## How to use it

- You can replace the `public` directory with any other html+js code that needs signaling server and use this as HTTP server and signaling server
- You can run it as signaling server and have another server serving your html/js/css that then connects to this signaling server from the client side.

## Installation

### Download latest binary

You can download [the latest version from releases](https://github.com/emad-elsaid/inbox/releases/latest) for your system/architecture

### Compile from source

- Have the go toolchain installed https://golang.org/dl/
- Clone the repository and compile and install the binary to $GOBIN
  ```
  git clone git@github.com:emad-elsaid/inbox.git
  cd inbox
  go install cmd/inbox.go
  ```

### Docker Image

- if you want to run it in HTTP mode
  ```
  docker run --rm -it -p 3000:3000 emadelsaid/inbox ./inbox --https=false
  ```
- You can use ssl-gen script to generate a self signed certificate, or if you already have a certificate
  ```
  docker run --rm -it -v /path/to/cert/directory:/cert -p 3000:3000 emadelsaid/inbox ./inbox --server-cert=/cert/server.crt --server-key=/cert/server.key
  ```

## Usage

```
  -bind string
        a bind for the http server (default "0.0.0.0:3000")
  -cleanup-interval int
        Interval in seconds between server cleaning up inboxes (default 1)
  -https
        Run server in HTTPS mode or HTTP (default true)
  -public string
        Directory path of static files to serve (default "public")
  -server-cert string
        HTTPS server certificate file (default "server.crt")
  -server-key string
        HTTPS server private key file (default "server.key")
```

## API Documentation

- Swagger documentation is under [/swagger.yml](/swagger.yml)
- You can show it live here https://petstore.swagger.io/ , the use the following
  URL in the top input field
  ```
  https://raw.githubusercontent.com/emad-elsaid/inbox/master/swagger.yml
  ```
