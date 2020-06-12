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
- it uses Ruby to run a web server on port 3000 that serves 2 pages `/send` for
  the camera machine and `/receive` for the receiver machine that wants to
  display the camera.
- It uses WebRTC to start a webRTC connection
- The local server has 4 other routes for the sender and receiver to signal each
  other the webRTC offer and answer.


## State

This is working but it's still too buggy and flacky, so it's in a state were it
needs a bit of work to handle some cases like reconnecting if the connection
failed.

## Help needed.

This project needs a little help to be more stable, the code isn't large at all,
the sinatra server is one file, there are 2 javascript files for sender and
receiver pages and another 2 files for signaling and webrtc common code.


## How to run

- Clone the code
- make sure you have ruby installed
- install gems `bundle install`
- run the server `./server`
- open `https://your-ip-address:3000/send` on the camera machine
- open `https://your-ip-address:3000/receive` on the receiver machine
- choose the camera from the list on the sender and press `start` button
- the receiver should display the camera shortly after

## Problems

- if the camera didn't appear on the receiver screen in 15 seconds try
  refreshing the receiver or both receiver and sender
- The webRTC offer isn't enough when it's sent for the first time so I have to
  keep sending new offer every second until the connection is established.
- the video sometimes doesn't play on the receiver until you interact with the
  page (google chrome policy)
