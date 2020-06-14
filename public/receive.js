async function RTCTrack(event) {
  console.log('RTCTrack', event);

  var remoteStream = new MediaStream();
  remoteStream.addTrack(event.track);

  var preview = document.getElementById('preview');
  preview.srcObject = remoteStream;
}

peer = new Peer(new SignalingChannel('receiver'));
peer.connection.addEventListener('track', RTCTrack);
