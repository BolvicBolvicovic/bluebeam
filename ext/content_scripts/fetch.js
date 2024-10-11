function analyze(username, sessionKey) {
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
    username: username,
    sessionKey: sessionKey,
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
      const jsonData = JSON.parse(data);
      browser.runtime.sendMessage({ type: 'analyzeResponse', data: jsonData });
  })
  .catch(error => console.error('Error sending data:', error));
}

async function login(message) {
  let username = message.username;
  let sessionKey;
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
      sessionKey = jsonData.session_key;

  })
  .catch(error => console.error('Error sending data:', error));
  return [username, sessionKey];
}

function isConnected(sessionKey) {
      let bool = sessionKey ? true : false;
      browser.runtime.sendMessage({ isConnected: bool });
}

function register(message) {
  let body = JSON.stringify({ 
      username: message.username,
      password: message.password,
  })

  fetch('https://localhost/register_account', {
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

function settingsPage(username, sessionKey) {
  const url = `https://localhost/settings?username=${encodeURIComponent(username)}&sessionkey=${encodeURIComponent(sessionKey)}`;
  window.open(url);
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
  let username;
  let sessionKey;
  browser.runtime.onMessage.addListener(async (message) => {
    if (message.type === "analyze") {
      analyze(username, sessionKey);
    } else if (message.type === "login") {
      values = await login(message);
      username = values[0];
      sessionKey = values[1];
    } else if (message.type === "isConnected") {
      isConnected(sessionKey);
    } else if (message.type === "register") {
      register(message);
    } else if (message.type === "settingsPage") {
      settingsPage(username, sessionKey);
    }
  });
})();
