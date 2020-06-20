class SignalingChannel extends EventTarget {
  constructor(opts) {
    super();
    this.from = opts['from'];
    this.to = opts['to'];
    this.password = opts['password'];
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

    const response = await fetch(`/inbox?from=${this.from}&to=${this.to}&password=${this.password}`, {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(data)
    });
  }

  async receive() {
    const response = await fetch(`/inbox?to=${this.from}&password=${this.password}`, {
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
