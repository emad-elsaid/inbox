class Peer extends EventTarget{
  constructor(channel) {
    super();
    this.connection = new RTCPeerConnection();
    this.connection.addEventListener('connectionstatechange', e => this.stateChanged(e));
    this.connection.addEventListener('icecandidate', e => this.sendIceCandidate(e));
    this.channel = channel;
    this.channel.addEventListener('offer', e => this.answer(e.detail));
    this.channel.addEventListener('answer', e => this.acceptAnswer(e.detail));
    this.channel.addEventListener('icecandidate', e => this.receiveIceCandidate(e.detail));
  }

  get status() {
    return this.connection.connectionState;
  }

  stateChanged(event) {
    console.log("Connection changed", this.connection.connectionState);
    this.dispatchEvent(new CustomEvent(this.connection.connectionState));
  }

  addStream(stream) {
    stream.getTracks().forEach(t => this.connection.addTrack(t));
  }

  async connect() {
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

  sendIceCandidate(event) {
    if( event.candidate == null ) return;

    this.channel.send({ type: 'icecandidate', candidate: event.candidate.toJSON() });
  }

  receiveIceCandidate(event) {
    if( event.candidate == null ) return;

    var candidate = new RTCIceCandidate(event.candidate);
    this.connection.addIceCandidate(candidate);
  }
}
