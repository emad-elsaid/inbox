# Tutorial

## Installation

* Download the latest version from github [the latest version from releases](https://github.com/emad-elsaid/inbox/releases/latest)
* Extract the zip file to current directory
* You should see the main binary `inbox` and examples directory `public`

## Generating self signed certificate for development

WebRTC is not allowed in browsers unless your connection is encrypted, inbox
runs in HTTPS mode by default given that you have certificate files in the
current directory.

Lets generate certificate files.
```
openssl genrsa -des3 -passout pass:secretpassword -out server.pass.key 2048
openssl rsa -passin pass:secretpassword -in server.pass.key -out server.key
openssl req -new -key server.key -out server.csr
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
```

When you're asked about any information no need to enter any of it, also don't
add a password/secret phrase to the certificate, remember this is for you on
your development machine.

## Running the server

* Start the server in the directory that contains the certificate files and the
  `public` directory
* The server will pickup the certificate files as long as it's named
  `server.key` and `server.crt`
* Also it will serve all files in the `public` directory
* On linux you can run it by executing `./inbox` it's this simple

## Accessing the server

* The server listen on port `3000`
* In your browser open `https://localhost:3000` make sure it's `https` not
  'http'
* You should see a list of files in `public` directory now `send.html` and
  `receive.html` should allow you to share the camera between two tabs give it a
  try, if you terminated the server after starting the connection it will
  continue to work as it doesn't need the server any more.

## Writing you WebRTC application from scratch

* The server serves all static files under `public` but you can pass `-public`
  with the directory path you want to serve instead.
* The server has CORS turned off by default you can turn it on by passing
  `-cors=true` so in case you want to split the signaling server (inbox) from
  your assets server this will be useful.
* Inbox has long polling turned on by default so when a user asks Inbox about a
  message it will wait until a message is received from another peer then it
  will respond with this message, you can turn off this feature to respond
  directly instead of waiting by passing `-long-polling=false`

## Signal Javascript class

To use inbox in your javascript application you can write a small class that
asks the server for new messages or send messages to another peer with his ID.

The following snippet will send a message to a peer by his user name:

```js
async send(server, from, password, to, data) {
  const response = await fetch(`${server}/inbox?to=${to}`, {
    method: 'POST',
    cache: 'no-cache',
    headers: new Headers({ 'Authorization': 'Basic ' + window.btoa(from + ":" + password) }),
    body: JSON.stringify(data)
  });
}
```

Inbox will get the message and will create an inbox for the `from` user with
`password` provided if the inbox doesn't exist.

To send a message to another user this user has to exist on the server, this is
why every user should first start by asking the server about any new messages,
then send a message to another user if he wants to.

The following function should allow you to get the latest message from the
server

```js
async receive(server, from, password) {
  const response = await fetch("${server}/inbox", {
    method: 'GET',
    cache: 'no-cache',
    headers: new Headers({ 'Authorization': 'Basic ' + window.btoa(from + ":" + password) }),
  });

  try {
    var data = await response.json();
    return data;
  } catch(e) {
    return null;
  }
}
```

This will get the latest message (the server will not respond until there is a
message, if inbox is empty it will wait until it gets a message from the other
end) the server will create an inbox for the user with the provided username and
password

You'll need to keep polling the server for new message as long as you expect new
messages from peers, if you are connected to all expected peers then no need to
ask the server anymore.

There is no technical limitation that stops you from continueing to poll the
server even after connecting to your peers, but it will consume server resources
so it's a good practice to be conservative with your requests to allow the
server to serve as many users as possible.

Your javascript code now can use these two functions to send and receive message
from other peers in conjunction with `RTCPeerConnection` javascript class to
generate WebRTC offers/answers and send it to the other peer until the
connection is established, you can check `/public/webrtc.js` for an example of
how to do that.
