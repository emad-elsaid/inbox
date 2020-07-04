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
  -cors
        Allow CORS
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

- Swagger documentation is under [/docs/swagger.yml](/docs/swagger.yml)
- You can show it live here https://petstore.swagger.io/ , the use the following
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
- Every second the server removes inboxes and messages exceeded timeouts

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

## Benchmarks

On a machine with the following specifications
```
$ lscpu
Architecture:                    x86_64
CPU op-mode(s):                  32-bit, 64-bit
Byte Order:                      Little Endian
Address sizes:                   39 bits physical, 48 bits virtual
CPU(s):                          4
On-line CPU(s) list:             0-3
Thread(s) per core:              2
Core(s) per socket:              2
Socket(s):                       1
NUMA node(s):                    1
Vendor ID:                       GenuineIntel
CPU family:                      6
Model:                           142
Model name:                      Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
Stepping:                        9
CPU MHz:                         2108.897
CPU max MHz:                     3900.0000
CPU min MHz:                     400.0000
BogoMIPS:                        5802.42
Virtualization:                  VT-x
L1d cache:                       64 KiB
L1i cache:                       64 KiB
L2 cache:                        512 KiB
L3 cache:                        4 MiB
NUMA node0 CPU(s):               0-3
Vulnerability Itlb multihit:     KVM: Mitigation: Split huge pages
Vulnerability L1tf:              Mitigation; PTE Inversion; VMX conditional cache flushes, SMT vulnerable
Vulnerability Mds:               Mitigation; Clear CPU buffers; SMT vulnerable
Vulnerability Meltdown:          Mitigation; PTI
Vulnerability Spec store bypass: Mitigation; Speculative Store Bypass disabled via prctl and seccomp
Vulnerability Spectre v1:        Mitigation; usercopy/swapgs barriers and __user pointer sanitization
Vulnerability Spectre v2:        Mitigation; Full generic retpoline, IBPB conditional, IBRS_FW, STIBP conditional, RSB filling
Vulnerability Tsx async abort:   Mitigation; Clear CPU buffers; SMT vulnerable
Flags:                           fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss
                                  ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonsto
                                 p_tsc cpuid aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid
                                  sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpu
                                 id_fault epb invpcid_single pti ssbd ibrs ibpb stibp tpr_shadow vnmi flexpriority ept vpid ept_ad fsgsbase ts
                                 c_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx rdseed adx smap clflushopt intel_pt xsaveopt xsavec xge
                                 tbv1 xsaves dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp md_clear flush_l1d```
```

Go benchmark command for 1 second produces the following results

```
goos: linux
goarch: amd64
pkg: inbox
BenchmarkInboxPut-4              3469288               308 ns/op
BenchmarkInboxPutThenGet-4       9196538               119 ns/op
BenchmarkServerGet-4             2177724               577 ns/op
BenchmarkServerPost-4            1626523               811 ns/op
PASS
ok      inbox   6.675s
```

## Contribute

Expected contribution flow will be as follows:

* Read an understand the code
* Make some changes related to your idea
* Open a PR to validate you're in the right way, describe what you're trying to do
* Continue working on your changes until it's fully implemented
* I'll merge it and release a new version

## License

MIT License (c) Emad Elsaid
