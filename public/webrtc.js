const peerConnection = new RTCPeerConnection();
peerConnection.addEventListener('connectionstatechange', RTCConnectionChanged);

function RTCConnectionChanged(event) {
  console.log("Connection Status changed", peerConnection.connectionState);

  if( peerConnection.connectionState == 'failed' ) {
    document.location.reload();
  }
}
