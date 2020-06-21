WebRTC Camera
=============

## Purpose

A webserver that you can run locally, you can use it to stream your camera from
one machine to the other over your network, it works from phone to PC, My
intention was to use my phone camera as a high resolution wide camera len for my Twitch
streams with OBS browser source instead of my webcamera.

## How is it working?

- the repository incudes a script to generate SSL self signed certificate
  `ssl-gen` as it's needed to run the server and use webRTC
- it uses Go to run a web server on port 3000 that serves 2 pages `/send.html` for
  the camera machine and `/receive.html` for the receiver machine that wants to
  display the camera.
- It uses WebRTC to start a webRTC connection
- The local server has 1 other route for the sender and receiver to signal each
  other the webRTC offer and answer.

## How to run

- Clone the code
- make sure you have Go installed
- run the server `go run ./cmd/server.go`
- open `https://your-ip-address:3000/send.html` on the camera machine
- open `https://your-ip-address:3000/receive.html` on the receiver machine
- choose the camera from the list on the sender and press `start` button
- the receiver should display the camera shortly after

## Problems

- Doesn't work on firefox
- the video sometimes doesn't play on the receiver until you interact with the
  page (google chrome policy) if you want to disable it you can run the receiver
  page as a google app with this policy disabled
  ```
  google-chrome-stable --app=https://server-ip-address:3000/receive --enable-features="PreloadMediaEngagementData,AutoplayIgnoreWebAudio,MediaEngagementBypassAutoplayPolicies"
  ```
