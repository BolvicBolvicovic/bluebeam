function reportError(error) {
    console.error(`Error caught: ${error}`);
}

function registerLoginButton() {
  document.getElementById("loginButton").addEventListener("click", (e)=> {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
      document.getElementById("consoleMessage").innerHTML = 'Please provide both username and password.';
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

function registerRegisterButton() {
  document.getElementById("registerButton").addEventListener("click", (e)=> {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
      document.getElementById("consoleMessage").innerHTML = 'Please provide both username and password.';
      return;
    }

    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage( tabs[0].id, {
          type: "register",
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
        document.getElementById("consoleMessage").innerHTML = 'Data sent, analysing...';
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

function authenticate() {
  const clientID = "726518157620-ue8o67ep33k2cr2k6ra8uqlpneo6uodu.apps.googleusercontent.com";
  const redirectURI = browser.identity.getRedirectURL();
  const authURL = `https://accounts.google.com/o/oauth2/auth?client_id=${clientId}&response_type=token&redirect_uri=${encodeURIComponent(redirectUri)}&scope=https://www.googleapis.com/auth/spreadsheets`;
  return browser.identity.lauchWebAuthFlow({
    interactive: true,
    url: authURL
  }).then(responseURL => {
    const params = new URL(responseURL).hash.substring(1);
    return URLSearchParams(params).get("access_token");
  });
}

function createSpreadSheet(accessToken) {
  const url = 'https://sheets.googleapis.com/v4/spreadsheets';
  const body = {
    properties: {
      // TODO: Add the name of the website we are on
      title: "data"
    }
  };
  return fetch(url, {
    method: 'POST',
    headers: {
      "Authorization": `Bearer ${accessToken}`,
      "Content-Type": "application/json"
    },
    body: JSON.stringify(body)
  })
  .then(response => response.json())
  .then(data => {
    return data.spreadsheetId;
  })
  .catch(e => console.error("Error creating spreadsheet:", e));
}

function updateSpreadsheet(accessToken, spreadsheetId, range, values) {
    const url = `https://sheets.googleapis.com/v4/spreadsheets/${spreadsheetId}/values/${range}:append?valueInputOption=RAW`;
    const body = JSON.stringify({
        values: values
    });

    return fetch(url, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json'
        },
        body: body
    })
    .then(response => response.json())
    .then(data => {
        console.log('Spreadsheet updated:', data);
    })
    .catch(error => console.error('Error updating spreadsheet:', error));
}

function convertTo2DArray(jsonArray) {
    if (jsonArray.length === 0) return [];
    
    const headers = Object.keys(jsonArray[0]);
    const rows = jsonArray.map(obj => headers.map(header => obj[header]));
    
    return [headers, ...rows];
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
    authenticate().then(accessToken => {
      createSpreadSheet(accessToken).then(spreadsheetId => {
        const range = 'Sheet1!A1';
        const values = convertJSONTo2DArray(data);
        updateSpreadsheet(accessToken, spreadsheetId, range, values);
      });
    });
  });
}

function messageListener() {
  browser.runtime.onMessage.addListener((message) => {
    if (message.type === "loginResponse") {
      if (message.data.error) {
        document.getElementById("consoleMessage").innerHTML = message.data.error;
        return;
      }
      document.getElementById("consoleMessage").innerHTML = "";
      document.getElementById("login").style.display = "none";
      document.getElementById("scrape").style.display = "block";
    } else if (message.type === "registerResponse") {
      document.getElementById("consoleMessage").innerHTML = (message.data.error != undefined) ? message.data.error : message.data.message;
    } else if (message.type === "analyzeResponse") {
      if (message.data.error != undefined) {
        document.getElementById("consoleMessage").innerHTML =  message.data.error;
      } else {
        buildDataFiles(message.data.message);
        document.getElementById("getOutput").style.display = "block";
      }
    } else if (message.isConnected === true) {
      document.getElementById("login").style.display = "none";
    } else if (message.isConnected === false) {
      document.getElementById("scrape").style.display = "none";
    } else if (message.error) {
      document.getElementById("consoleMessage").innerHTML = message.error;
    }
  });
}

async function handler() {
  await browser.tabs.executeScript({ file: "/content_scripts/fetch.js" })
  document.getElementById("getOutput").style.display = "none";
  messageListener();
  requestIsConnected();
  registerLoginButton();
  registerRegisterButton();
  registerScrapeButton();
  registerSettings();
}

handler()
