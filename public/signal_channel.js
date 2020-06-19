class SignalingChannel extends EventTarget {
  constructor(role) {
    super();
    this.role = role;
    this.connected = false;
    this.connect();
  }

  connect() {
    if ( !this.connected ) {
      this.connected = true;
      this.poll();
    }
  }

  disconnect() {
    this.connected = false;
  }

  async send(data) {
    console.log('Sending', data);

    const response = await fetch(`/from/${this.role}`, {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(data)
    });
  }

  async receive() {
    const response = await fetch(`/inbox/${this.role}`, {
      method: 'GET',
      cache: 'no-cache'
    });

    try {
      var data = await response.json();
      console.log('Received', data);
      return data;
    } catch(e) {
      return null;
    }
  }

  async poll() {
    var message = await this.receive();
    if ( message != null ) {
      this.dispatchEvent(new CustomEvent(message.type || 'message', { detail: message }));
    }

    if ( this.connected ) setTimeout(this.poll.bind(this), 1000);
  }
}
