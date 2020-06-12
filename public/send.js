async function videoDevices() {
  var all =  await navigator.mediaDevices.enumerateDevices();
  var video = all.filter(d => d.kind == 'videoinput');
  return video;
}

async function videoStream(deviceId) {
  return await navigator.mediaDevices.getUserMedia({
    video: {
      width: { ideal: 4096 },
      height: { ideal: 4208 },
      deviceId: { exact: deviceId }
    }
  });
}

async function listDevices() {
  var dom = document.getElementById('devices');
  var devices = await videoDevices();
  var options = devices.map(videoDevice => {
    return `<option value="${videoDevice.deviceId}">${videoDevice.label}</option>`;
  });
  dom.innerHTML = options.join('');
}

async function start() {
  var deviceSelect = document.getElementById('devices');
  var deviceId = deviceSelect.value;
  var device = await videoStream(deviceId);
  var video = document.getElementById('preview');
  video.srcObject = device;

  peerConnection.addStream(new MediaStream());
  makeCall();
}
document.getElementById('start').addEventListener('click', start);

async function makeCall() {
  const offer = await peerConnection.createOffer();
  await peerConnection.setLocalDescription(offer);
  signalingChannel.send(offer.toJSON());

  // Send offer many time while ICE information is collected
  // until the connection is successfull
  if ( peerConnection.connectionState != 'connected' ) {
    setTimeout(makeCall, 2000);
  }
}

async function answerReceived(answer) {
  const remoteDesc = new RTCSessionDescription(answer);
  await peerConnection.setRemoteDescription(remoteDesc);

  var video = document.getElementById('preview');
  video.srcObject.getTracks().forEach(track => {
    peerConnection.addTrack(track);
  });
}

function messagesHandler(message) {
  console.log(message);

  switch(message.type) {
  case 'answer':
    answerReceived(message);
    break;
  default:
    console.error('Unknow message type', message);
  }
}

const signalingChannel = new SignalingChannel('sender');
signalingChannel.addEventListener(messagesHandler);

(async function() {
  await navigator.mediaDevices.getUserMedia({video: true});
  listDevices();
})();
