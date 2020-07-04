# The Concept

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
- A peer can use another peer username to send a message to his inbox
- The peer inbox will expire after a period of time (1 minute by default) of not
  asking for any message
- The message has a timeout and will be deleted after this timeout (1 minute by default)
- So peers has to keep asking the server for new messages with short delays that
  doesn't exceed the timeout until they got enough information to connect to
  each other

# Usecase
- Assuming 2 peers (Alice and Bob) want to connect with WebRTC
- **Alice** need to choose an identifier `alice-uuid` and pass it to **Bob** in any
  other medium (Chat or write it on a paper or pre share it)
- Alice uses `alice-uuid` as username to create her inbox and wait for messages
  from any peer **Bob** in our case
- **Bob** will create an inbox with any username `bob-uuid` and sends WebRTC offer to
  initiate connect with the pre shared username `bob-uuid`.
- **Alice** will ask the server for new messages with her username `alice-uuid`
- The server responds with **Bob** WebRTC offer message in reponse body and
  `X-From` header with `bob-uuid` as value
- **Alice** will send WebRTC answer to `bob-uuid`
- **Bob** Asks the server for new messages
- The server responds with **Alice** webRTC answer message and `X-From` header with `alice-uuid`
- **Bob** sends **Alice** ICE candidates information in a message each time he's
  aware of new candidate
- **Alice** will receive ICE candidates messages and sends **Bob** candidates messages
- **Alice** and **Bob** keep sending each other messages through the server
  until they have enough information to connect to each other
- **Alice** and **Bob** are now connected to each other directly with WebRTC
- The server will delete both inboxes after 1 minute of inactivity
