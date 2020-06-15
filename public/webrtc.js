class Peer {
  constructor(channel) {
    this.connection = new RTCPeerConnection();
    this.connection.addEventListener('connectionstatechange', (e) => this.stateChanged(e));
    this.channel = channel;
    this.channel.addEventListener('offer', (e) => this.answer(e.detail));
    this.channel.addEventListener('answer', (e) => this.acceptAnswer(e.detail));
    this.channel.startPolling();
  }

  stateChanged(event) {
    console.log("Connection changed", this.connection.connectionState);

    switch( this.connection.connectionState ) {
    case 'connected':
      this.channel.stopPolling();
      break;
    case 'failed':
    case 'disconnected':
      document.location.reload();
      break;
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
}
