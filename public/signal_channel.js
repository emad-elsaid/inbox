class SignalingChannel extends EventTarget {
  constructor(opts) {
    super();
    this.from = opts['from'];
    this.to = opts['to'];
    this.password = opts['password'];
    this.connected = false;
    this.headers = new Headers({ 'Authorization': 'Basic ' + window.btoa(this.from + ":" + this.password) });
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

    const response = await fetch(`/inbox?to=${this.to}`, {
      method: 'POST',
      cache: 'no-cache',
      headers: this.headers,
      body: JSON.stringify(data)
    });
  }

  async receive() {
    const response = await fetch("/inbox", {
      method: 'GET',
      cache: 'no-cache',
      headers: this.headers
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
    if ( message != null ) this.dispatchEvent(new CustomEvent(message.type || 'message', { detail: message }));
    if ( this.connected ) setTimeout(this.poll.bind(this), 1000);
  }
}
