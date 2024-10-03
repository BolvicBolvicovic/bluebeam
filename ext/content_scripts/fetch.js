function analyze(message) {
  let links = Array.from(document.querySelectorAll('a')).map(a => a.href);
  let buttons = Array.from(document.querySelectorAll('button')).map(b => ({
    text: b.innerText,
    onclick: b.onclick ? b.onclick.toString() : null
  }));
  let pageHtml = document.documentElement.outerHTML;
  let body = JSON.stringify({ 
      username: message.username,
      sessionkey: message.sessionKey,
      links,
      buttons,
      pageHtml
  })

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
      browser.runtime.sendMessage({ type: 'analyzeResponse', data: jsonData });
  })
  .catch(error => console.error('Error sending data:', error));
}

function login(message) {
  let body = JSON.stringify({ 
      username: message.username,
      password: message.password,
  })

  fetch('https://localhost/login', {
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
  browser.runtime.onMessage.addListener((message) => {
    if (message.type === "analyze") {
      analyze(message);
    } else if (message.type === "login") {
      login(message);
    }

  });
})();
