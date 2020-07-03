📮 INBOX
=============

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/inbox)](https://goreportcard.com/report/github.com/emad-elsaid/inbox)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/inbox?status.svg)](https://godoc.org/github.com/emad-elsaid/inbox)
[![codecov](https://codecov.io/gh/emad-elsaid/inbox/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/inbox)
[![Join the chat at https://gitter.im/inbox-server/community](https://badges.gitter.im/inbox-server/community.svg)](https://gitter.im/inbox-server/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Inbox makes it easy to setup a WebRTC HTTPS signaling server

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

- The server acts as temporary mailbox for peers
- Peers use basic authentication (username, password) to get or send messages
- Username can be anything: random number, UUID, public key...etc
- Whenever a peer authenticate with username and password an inbox will be
  created for them if it doesn't exist
- If the username exists and the password is correct then the server will
  respond with the oldest message in the inbox and deletes it from it's memory,
  and will respond with header `X-From` with the peer username that sent this
  message.
- If the username exists and the password is incorrect an Unauthorized arror is
  returned
- Now the Inbox with this username is ready to receive messages from another
  peer.
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
- No thir party dependencies at all
- Stores all data in memory in one big memory structure
- Clears old data every second to remove inboxes and messages exceeded timeouts

## How to run the example

- Install [Go](https://golang.org/)
- Clone this repository `git clone git@github.com:emad-elsaid/inbox.git`
- Run the server `go run ./cmd/inbox.go`
- Open `https://your-ip-address:3000/send.html` on the camera machine
- Open `https://your-ip-address:3000/receive.html` on the receiver machine
- Choose the camera from the list on the sender and press `start` button
- The receiver should display the camera shortly after

## How to use it

- You can replace the `public` directory with any other html+js code that needs signaling server and use this as http server and signaling server
- You can run it as signaling server and have another server serving your html/js/css that then connects to this signaling server from the client side.

## installation

### download latest binary

You can download [the latest version from releases](https://github.com/emad-elsaid/inbox/releases/latest) for your system/architecture

### compile from source

- Have the [Go toolchain](https://golang.org/dl/) installed
- Clone the repository and compile and install the binary to $GOBIN
  ```
  git clone git@github.com:emad-elsaid/inbox.git
  cd inbox
  go install cmd/inbox.go
  ```

### docker image

- If you want to run it in http mode
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

## Benchmarks

On a machine with the following specifications
```
$ lscpu
Architecture:                    x86_64
CPU op-mode(s):                  32-bit, 64-bit
Byte Order:                      Little Endian
Address sizes:                   48 bits physical, 48 bits virtual
CPU(s):                          8
On-line CPU(s) list:             0-7
Thread(s) per core:              1
Core(s) per socket:              8
Socket(s):                       1
NUMA node(s):                    1
Vendor ID:                       AuthenticAMD
CPU family:                      23
Model:                           8
Model name:                      AMD Ryzen 7 2700X Eight-Core Processor
Stepping:                        2
CPU MHz:                         3693.050
BogoMIPS:                        7389.85
Hypervisor vendor:               KVM
Virtualization type:             full
L1d cache:                       256 KiB
L1i cache:                       512 KiB
L2 cache:                        4 MiB
L3 cache:                        16 MiB
NUMA node0 CPU(s):               0-7
Vulnerability Itlb multihit:     Not affected
Vulnerability L1tf:              Not affected
Vulnerability Mds:               Not affected
Vulnerability Meltdown:          Not affected
Vulnerability Spec store bypass: Mitigation; Speculative Store Bypass disabled via prctl and seccomp
Vulnerability Spectre v1:        Mitigation; usercopy/swapgs barriers and __user pointer sanitization
Vulnerability Spectre v2:        Mitigation; Full AMD retpoline, STIBP disabled, RSB filling
Vulnerability Srbds:             Not affected
Vulnerability Tsx async abort:   Not affected
Flags:                           fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt rdtscp lm constant_tsc rep_good
                                  nopl nonstop_tsc cpuid extd_apicid tsc_known_freq pni pclmulqdq ssse3 cx16 sse4_1 sse4_2 x2apic movbe popcnt aes xsave avx rdrand hypervisor lahf_lm cmp_legac
                                 y cr8_legacy abm sse4a misalignsse 3dnowprefetch ssbd vmmcall fsgsbase avx2 rdseed clflushopt arat npt lbrv svm_lock nrip_save tsc_scale vmcb_clean flushbyasid
                                  decodeassists pausefilter pfthreshold avic v_vmsave_vmload vgif
```

Go benchmark command for 1 second produces the following results

```
go test -bench . -benchtime=1s
goos: linux
goarch: amd64
pkg: inbox
BenchmarkInboxPut-8              1836072               666 ns/op
BenchmarkInboxPutThenGet-8       8294449               146 ns/op
BenchmarkServerGet-8             1928589               622 ns/op
BenchmarkServerPost-8            1408486               779 ns/op
PASS
ok      inbox   7.092s
```
