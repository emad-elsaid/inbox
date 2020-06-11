peerConnection.addEventListener('track', RTCTrack);

async function makeCall(offer) {
  peerConnection.setRemoteDescription(new RTCSessionDescription(offer));
  const answer = await peerConnection.createAnswer();
  peerConnection.setLocalDescription(answer);
  signalingChannel.send(answer.toJSON());
};

async function RTCTrack(event) {
  console.log('RTCTrack', event);

  var remoteStream = new MediaStream();
  remoteStream.addTrack(event.track);

  var preview = document.getElementById('preview');
  preview.srcObject = remoteStream;
  preview.play();
}

const signalingChannel = new SignalingChannel('receiver');
signalingChannel.addEventListener(messagesHandler);

function messagesHandler(message) {
  console.log(message);

  switch(message.type) {
  case 'offer':
    makeCall(message);
    break;
  default:
    console.error('Unknow message type', message);
  }
}
