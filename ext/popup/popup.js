let username;
let password;
let sessionKey;

function registerLoginButton() {
  alert("You can register with curl at https://localhost/register_account with a JSON username, password");
  document.getElementById("loginButton").addEventListener("click", (e)=> {
    username = document.getElementById('username').value;
    password = document.getElementById('password').value;

    if (!username || !sessionKey) {
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

function registerScrapeButton(scrapeButton) {
  document.getElementById('scrapeButton').addEventListener("click", (e) => {
    function reportError(error) {
      console.error(`Error caught: ${error}`);
    }
    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage( tabs[0].id, {
          type: "analyze",
          username: username,
          sessionKey: sessionKey,
        });
      })
      .then(() => {
        alert('Data sent to server');
      })
      .catch(reportError);
  });
}

function messageListener() {
  browser.runtime.onMessage.addListener((message) => {
    if (message.type === "loginResponse") {
      sessionKey = message.data.sessionKey;
      document.getElementById("login").style.display = "none";
      document.getElementById("scrape").style.display = "block";
    } else if (message.type === "analyzeResponse") {
      alert(message.data)
    } else if (message.error) {
      alert(message.error)
    }
  });
}

function handler() {
  document.getElementById("scrape").style.display = "none";
  registerLoginButton();
  registerScrapeButton();
}

browser.tabs
  .executeScript({ file: "/content_scripts/fetch.js" })
  .then(handler)
