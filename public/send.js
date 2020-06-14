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

  updateCapabilitiesForm();
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

function updateCapabilitiesForm() {
  var form = document.getElementById('capabilities');
  var preview = document.getElementById('preview');
  var track = preview.srcObject.getVideoTracks()[0];
  var capabilities = track.getCapabilities();
  var settings = track.getSettings();

  var inputs = [];
  for( var capability in capabilities ) {
    var value = capabilities[capability];
    inputs.push("<div>");

    if( Array.isArray(value) ) {
      inputs.push(`<label for="${capability}">${capability}</label>`);
      inputs.push(`<select name="${capability}">`)
      value.forEach( val => {
        if( settings[capability] === val ) {
          inputs.push(`<option value="${val}" selected>${val}</option>`)
        } else {
          inputs.push(`<option value="${val}">${val}</option>`)
        }
      });
      inputs.push(`</select>`)

    }else if ( typeof value === 'object' ) {
      inputs.push(`<label for="${capability}">${capability}</label>`);
      inputs.push(`<input name="${capability}" type="range" min="${value.min}" max="${value.max}" value="${settings[capability]}">`)
    }

    inputs.push("</div>");
  }
  inputs.push(`<button id="updateCapabilities" type="submit">Update</button>`)

  form.innerHTML = inputs.join('');
}

const signalingChannel = new SignalingChannel('sender');
signalingChannel.addEventListener(messagesHandler);

(async function() {
  var media = await navigator.mediaDevices.getUserMedia({video: true});
  media.getVideoTracks()[0].stop();
  listDevices();
})();

document.getElementById("capabilities").addEventListener('submit', async function(event) {
  event.preventDefault();
  var constraints = {}
  var elements = this.elements
  for( var i = 0; i < elements.length; i++) {
    var element = elements[i];
    if ( element.name !== '') {
      constraints[element.name] = element.value;
    }
  }
  console.log("Applying constraints", constraints);

  var preview = document.getElementById('preview');
  var track = preview.srcObject.getVideoTracks()[0];
  await track.applyConstraints(constraints);

  updateCapabilitiesForm();
})
