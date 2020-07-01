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

  peer.addStream(device);
  peer.connect();
}
document.getElementById('start').addEventListener('click', start);

(async function() {
  var media = await navigator.mediaDevices.getUserMedia({video: true});
  media.getVideoTracks()[0].stop();
  listDevices();
})();

signaling = new SignalingChannel({
  from: 'sender',
  to: 'receiver',
  password: 'secretpassword'
});
peer = new Peer(signaling);

peer.addEventListener('connected', (e) => signaling.disconnect());
peer.addEventListener('failed', (e) => document.location.reload());
peer.addEventListener('disconnected', (e) => document.location.reload());
