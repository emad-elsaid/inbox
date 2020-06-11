const stunServers = [
  {
    'urls': [
      'stun:stun.l.google.com:19302',
      'stun:global.stun.twilio.com:3478'
    ]
  }
];

const configuration = {
  iceServers: stunServers,
  sdpSemantics: 'unified-plan'
};

const peerConnection = new RTCPeerConnection(configuration);
peerConnection.addEventListener('connectionstatechange', RTCConnectionChanged);
function RTCConnectionChanged(event) {
  console.log("Connection Status changed", peerConnection.connectionState);
}
