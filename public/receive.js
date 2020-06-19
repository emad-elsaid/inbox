async function RTCTrack(event) {
  console.log('RTCTrack', event);

  var remoteStream = new MediaStream();
  remoteStream.addTrack(event.track);

  var preview = document.getElementById('preview');
  preview.srcObject = remoteStream;
}

signaling = new SignalingChannel('receiver');
peer = new Peer(signaling);
peer.connection.addEventListener('track', RTCTrack);

peer.addEventListener('connected', (e) => signaling.disconnect());
peer.addEventListener('failed', (e) => document.location.reload());
peer.addEventListener('disconnected', (e) => document.location.reload());
