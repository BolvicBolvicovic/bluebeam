function reportError(error) {
    console.error(`Error caught: ${error}`);
}


function registerLoginButton() {
  document.getElementById("loginButton").addEventListener("click", (e)=> {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
      alert('Please provide both username and password.');
      return;
    }

    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage( tabs[0].id, {
          type: "login",
          username: username,
          password: password,
        });
      })
      .catch(reportError);
  })
}

function registerScrapeButton() {
  document.getElementById('scrapeButton').addEventListener("click", (e) => {
    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage( tabs[0].id, {
          type: "analyze",
        });
      })
      .then(() => {
        alert('Data sent to server');
      })
      .catch(reportError);
  });
}

function requestIsConnected() {
  browser.tabs
    .query({ active: true, currentWindow: true })
    .then((tabs) => {
      browser.tabs.sendMessage( tabs[0].id, {
        type: "isConnected",
      });
    })
    .catch(reportError);
}

function registerSettings() {
  document.getElementById("settingsPage").addEventListener('click', (e) => {
      browser.tabs
        .query({ active: true, currentWindow: true })
        .then((tabs) => {
          browser.tabs.sendMessage(tabs[0].id, {
            type: "settingsPage",
          });
        })
        .catch(reportError);
  });
}

function messageListener() {
  browser.runtime.onMessage.addListener((message) => {
    if (message.type === "loginResponse") {
      document.getElementById("login").style.display = "none";
      document.getElementById("scrape").style.display = "block";
    } else if (message.type === "analyzeResponse") {
      console.log(message.data)
    } else if (message.isConnected === true) {
      document.getElementById("login").style.display = "none";
    } else if (message.isConnected === false) {
      document.getElementById("scrape").style.display = "none";
    } else if (message.error) {
      console.log(message.error)
    }
  });
}

async function handler() {
  await browser.tabs.executeScript({ file: "/content_scripts/fetch.js" })
  messageListener();
  requestIsConnected();
  registerLoginButton();
  registerScrapeButton();
  registerSettings();
}

handler()
