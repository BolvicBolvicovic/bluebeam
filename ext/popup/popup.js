const login = document.getElementById("login");
const register = document.getElementById("register");
const consoleMessage = document.getElementById("consoleMessage");


function reportError(error) {
    console.error(`Error caught: ${error}`);
}

function registerLoginButtons() {
  document.getElementById("loginButton").addEventListener("click", (e)=> {
    const username = document.getElementById('lUsername').value;
    const password = document.getElementById('lPassword').value;

    if (!username || !password) {
      consoleMessage.innerHTML = 'Please provide both username and password.';
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
  });
  document.getElementById("lRegisterButton").addEventListener("click", () => {
    consoleMessage.innerHTML = "";
    document.getElementById("login").style.display = "none";
    document.getElementById("register").style.display = "block";
  });
}

function registerRegisterButtons() {
  document.getElementById("rRegisterButton").addEventListener("click", (e)=> {
    const username = document.getElementById('rUsername').value;
    const password = document.getElementById('rPassword').value;
    const password2 = document.getElementById('rPassword2').value;
    const email = document.getElementById("email").value;

    if (!username || !password || !email) {
      consoleMessage.innerHTML = 'Please fill all fields.';
      return;
    }

    if (password != password2) {
      consoleMessage.innerHTML = 'The two passwords are different.';
      return;
    }

    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage( tabs[0].id, {
          type: "register",
          username: username,
          password: password,
          email: email
        });
      })
      .catch(reportError);
  });
  document.getElementById("backButton").addEventListener("click", () => {
    consoleMessage.innerHTML = "";
    document.getElementById("register").style.display = "none";
    document.getElementById("login").style.display = "block";
  });
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
        consoleMessage.innerHTML = 'Data sent, analysing...';
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

function buildDataFiles(data) {
  // JSON
  const dataStr = JSON.stringify(data, null, 2);
  const blob = new Blob([dataStr], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  document.getElementById("getJSON").addEventListener("click", () => {
    window.open(url, "_blank");
  });

  //Google Sheet
  document.getElementById("getGoogleSpreadsheet").addEventListener("click", () => {
    consoleMessage.innerHTML = "Data sent to Google, a new page will open soon...";
    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage(tabs[0].id, {
          type: "outputGoogleSpreadsheet",
          data: data
        });
      })
      .catch(reportError);
  });
}

function messageListener() {
  browser.runtime.onMessage.addListener((message) => {
    if (message.type === "loginResponse") {
      if (message.data.error) {
        consoleMessage.innerHTML = message.data.error;
        return;
      }
      consoleMessage.innerHTML = "";
      document.getElementById("login").style.display = "none";
      document.getElementById("scrape").style.display = "block";
    } else if (message.type === "registerResponse") {
      consoleMessage.innerHTML = (message.data.error != undefined) ? message.data.error : message.data.message;
    } else if (message.type === "analyzeResponse") {
      if (message.data.error != undefined) {
        consoleMessage.innerHTML =  message.data.error;
      } else {
        buildDataFiles(message.data.message);
        document.getElementById("getOutput").style.display = "block";
        consoleMessage.innerHTML = "";
      }
    } else if (message.isConnected === true) {
      document.getElementById("login").style.display = "none";
    } else if (message.isConnected === false) {
      document.getElementById("scrape").style.display = "none";
    } else if (message.error) {
      consoleMessage.innerHTML = message.error;
    }
  });
}

async function handler() {
  await browser.tabs.executeScript({file: "/content_scripts/fetch.js"});
  document.getElementById("getOutput").style.display = "none";
  document.getElementById("register").style.display = "none";
  messageListener();
  requestIsConnected();
  registerLoginButtons();
  registerRegisterButtons();
  registerScrapeButton();
  registerSettings();
}

handler()
