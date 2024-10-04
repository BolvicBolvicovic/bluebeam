
function analyze(username, sessionKey) {
  let links = Array.from(document.querySelectorAll('a')).map(a => a.href);
  let buttons = Array.from(document.querySelectorAll('button')).map(b => ({
    text: b.innerText,
    onclick: b.onclick ? b.onclick.toString() : null
  }));
  let pageHtml = document.documentElement.outerHTML;
  let body = JSON.stringify({ 
      username: username,
      sessionkey: sessionKey,
      links,
      buttons,
      pageHtml
  })
  console.log(sessionKey)

  // Send the scraped data along with the username and session key to your Go server
  fetch('https://localhost/analyze', {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: body
  })
  .then(response => response.text())
    .then(data => {
      const jsonData = JSON.parse(data);
      console.log(jsonData);
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
  console.log(username, sessionKey);
  return [username, sessionKey];
}

function isConnected(sessionKey) {
      let bool = sessionKey ? true : false;
      browser.runtime.sendMessage({ isConnected: bool });
}

function sendCriterias(username, sessionKey, body) {
  fetch('https://localhost/login', {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    username: username,
    sessionkey: sessionKey,
    body: body,
  })
  .then(response => response.text())
    .then(data => {
      const jsonData = JSON.parse(data);
      browser.runtime.sendMessage({ type: 'loginResponse', data: jsonData });
      sessionKey = jsonData.session_key;

  })
  .catch(error => console.error('Error sending data:', error));
  
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
    } else if (message.type === "criterias") {
      sendCriterias(message.body);
    }

  });
})();
