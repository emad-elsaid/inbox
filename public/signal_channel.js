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

  removeEventListener(listener) {
    delete this.listeners[this.listeners.indexOf(listener)];
    this.listeners = this.listeners.filter(e => e);
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
    if ( this.listeners.length > 0 ) {
      var message = await this.receive();
      if ( message != null ) {
        this.listeners.forEach( listener => listener(message));
      }
      setTimeout(this.poll.bind(this), 1000);
    }
  }
}
