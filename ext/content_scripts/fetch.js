function analyze() {
  let links = Array.from(document.querySelectorAll('a')).map(a => ({
    href: a.href,
    text: a.innerText
  }));

  let buttons = Array.from(document.querySelectorAll('button')).map(b => ({
    text: b.innerText,
    onclick: b.onclick ? b.onclick.toString() : null,
    id: b.id || null,
    classes: b.className || null
  }));

  let images = Array.from(document.querySelectorAll('img')).map(img => ({
    src: img.src,
    alt: img.alt || null,
    classes: img.className || null
  }));

  let formInputs = Array.from(document.querySelectorAll('input')).map(input => ({
    type: input.type,
    name: input.name || null,
    value: input.value || null
  }));

  let metaTags = Array.from(document.querySelectorAll('meta')).map(meta => ({
    name: meta.getAttribute('name') || meta.getAttribute('property') || null,
    content: meta.getAttribute('content') || null
  }));

  let headers = Array.from(document.querySelectorAll('h1, h2, h3, h4, h5, h6')).map(header => ({
    tag: header.tagName,
    text: header.innerText
  }));

  let bodyText = document.body.innerText;

  let dataPayload = JSON.stringify({ 
    links,
    buttons,
    images,
    formInputs,
    metaTags,
    headers,
    bodyText
  });

  fetch('https://localhost/analyze', {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: dataPayload
  })
  .then(response => response.text())
    .then(data => {
      const jsonData = data ? JSON.parse(data) : "";
      browser.runtime.sendMessage({ type: 'analyzeResponse', data: jsonData });
  })
  .catch(error => console.error('Error sending data:', error));
}

async function login(message) {
  let body = JSON.stringify({ 
      username: message.username,
      password: message.password,
  })

  await fetch('https://localhost/login', {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: body
  })
  .then(response => response.text())
    .then(data => {
      const jsonData = JSON.parse(data);
      browser.runtime.sendMessage({ type: 'loginResponse', data: jsonData });

  })
  .catch(error => console.error('Error sending data:', error));
}

function register(message) {
  let body = JSON.stringify({ 
      email: message.email
  });

  fetch('https://localhost/registerAccount', {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: body
  })
  .then(response => response.text())
    .then(data => {
      const jsonData = JSON.parse(data);
      browser.runtime.sendMessage({ type: 'registerResponse', data: jsonData });
  })
  .catch(error => console.error('Error sending data:', error));
}

function settingsPage() {
  const url = `https://localhost/settings`;
  window.open(url);
}

function pingServer() {
  fetch("https://localhost/ping")
  .then(response => response.json())
  .then(data => {
    if (!data.error) {
      browser.runtime.sendMessage({ isConnected: true });
    } else {
      browser.runtime.sendMessage({ isConnected: false });
    }
  })
}

function outputGoogleSpreadsheet(data) {
  let body = JSON.stringify({
    data: data
  });

  fetch('https://localhost/outputGoogleSpreadsheet', {
    method: "POST",
    mode: "cors",
    headers: { "Content-Type": "application/json" },
    body: body
  })
  .then(response => response.json())
  .then(data => {
    if (!data.error) {
      window.open(data.spreadsheetUrl);
    } else {
      console.error(data.error);
    }
  })
  .catch(e => console.error(e));
}

(() => {
  /**
   * Check and set a global guard variable.
   * If this content script is injected into the same page again,
   * it will do nothing next time.
   */
  if (window.hasRun) {
    return;
  }
  window.hasRun = true;
  browser.runtime.onMessage.addListener(async (message) => {
    if (message.type === "analyze") {
      analyze();
    } else if (message.type === "login") {
      login(message);
    } else if (message.type === "register") {
      register(message);
    } else if (message.type === "outputGoogleSpreadsheet") {
      outputGoogleSpreadsheet(message.data);
    } else if (message.type === "settingsPage") {
      settingsPage();
    } else if (message.type === "isConnected") {
      pingServer();
    }
  });
})();
