class SignalingChannel {
  constructor(role) {
    this.role = role;
    this.startedPolling = false;
    this.listeners = [];
  }

  addEventListener(listener) {
    this.listeners.push(listener);

    if (!this.startedPolling) {
      this.poll();
      this.startedPolling = true;
    }
  }

  async send(data) {
    const response = await fetch(`/from/${this.role}`, {
      method: 'POST',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });
  }

  async receive(data) {
    const response = await fetch(`/inbox/${this.role}`, {
      method: 'GET',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    try {
      return await response.json();
    } catch(e) {
      return null;
    }
  }


  async poll() {
    // when we're already connected we don't need to ask for more messages from the server
    if ( peerConnection.connectionState != 'connected' ) {
      var message = await this.receive();
      if ( message != null ) {
        this.listeners.forEach( listener => listener(message));
      }
      setTimeout(this.poll.bind(this), 1000);
    } else {
      setTimeout(this.poll.bind(this), 5000);
    }
  }
}
