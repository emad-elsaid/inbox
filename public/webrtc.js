class Peer {
  constructor(channel) {
    this.connection = new RTCPeerConnection();
    this.connection.addEventListener('connectionstatechange', (e) => this.stateChanged(e));
    this.channel = channel;
    this.channel.addEventListener(this.messagesHandler);
  }

  stateChanged(event) {
    console.log("Connection changed", this.connection.connectionState);

    if( this.connection.connectionState === 'failed' ) {
      document.location.reload();
    } else if( this.connection.connectionState == 'connected' ) {
      this.channel.removeEventListener(this.messagesHandler);
    }
  }

  get connectionState() {
    return this.connection.connectionState;
  }

  async offer() {
    const peerOffer = await this.connection.createOffer();
    await this.connection.setLocalDescription(peerOffer);
    this.channel.send(peerOffer.toJSON());
  }

  async answer(peerOffer) {
    this.connection.setRemoteDescription(new RTCSessionDescription(peerOffer));

    const answer = await this.connection.createAnswer();
    this.connection.setLocalDescription(answer);
    this.channel.send(answer.toJSON());
  }

  async acceptAnswer(peerAnswer) {
    const remoteDesc = new RTCSessionDescription(peerAnswer);
    await this.connection.setRemoteDescription(remoteDesc);
  }

  streamVideo(stream) {
    stream.getTracks().forEach(t => { this.connection.addTrack(t) });
  }

  messagesHandler = (function(message) {
    console.log(message);

    switch(message.type) {
    case 'offer':
      this.answer(message);
      break;
    case 'answer':
      this.acceptAnswer(message);
      break;
    default:
      console.error('Unknow message type', message);
    }
  }).bind(this)
}
