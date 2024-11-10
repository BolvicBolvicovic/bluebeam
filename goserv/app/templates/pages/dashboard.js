const urls = document.getElementById("urls");
const textArea = document.getElementById("textArea");
const inputFilesLinks = document.getElementById("inputFilesLinks");

const dropdownInputChoice = document.getElementById("dropdownInputChoice");
const formInputChoice = document.getElementById("formInputChoice");

const dropdownAIChoice = document.getElementById("dropdownAIChoice");
const formAIChoice = document.getElementById("formAIChoice");

const outputFilesLinks = document.getElementById("outputFilesLinks");
const popupOutput = document.getElementById("popupOutput");

const page = document.getElementById("page");

popupOutput.style.display = "none";

function handlePopup(data) {
  const spreadSheetButton = document.getElementById("spreadSheetButton");
  const JSONButton = document.getElementById("JSONButton");
  const messagePopup = document.getElementById("messagePopup");
  // JSON
  const dataStr = JSON.stringify(data, null, 2);
  const blob = new Blob([dataStr], { type: "application/json" });
  const url = URL.createObjectURL(blob);

  JSONButton.onclick = () => {
    popupOutput.style.display = "none";
    window.open(url, "_blank");
  };
  //Google Sheet
  spreadSheetButton.onclick = () => {
    messagePopup.innerHTML = "Data sent to Google, a new page will open soon...";
    fetch('/outputGoogleSpreadsheet', {
      method: "POST",
      mode: "cors",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({data: data})
    })
    .then(response => response.json())
    .then(data => {
      if (!data.error) {
        popupOutput.style.display = "none";
        window.open(data.spreadsheetUrl);
	load_output_files();
      } else {
        console.error(data.error);
      }
    })
    .catch(e => console.error(e));
  };
  
  popupOutput.style.display = "flex";
}

function remove_all_children(node) {
	while (node.firstChild) {
		node.removeChild(node.lastChild);
	}
}

function load_input_files() {
  fetch("/dashboard/inputFiles")
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        inputFilesLinks.innerHTML = data.error;
        dropdownInputChoice.disabled = true;              
      } else if (data.index === -1) {
        dropdownInputChoice.disabled = true;            
        return;
      } else {
        dropdownInputChoice.disabled = false;
        inputFilesLinks.innerHTML = "";
	remove_all_children(dropdownInputChoice);
	remove_all_children(inputFilesLinks);
        let i = 0;
        data.files.forEach(file => {
          const option = document.createElement("option")
          const fileLink = document.createElement("a");

          option.value = i;
          option.innerText = file.filename;
          if (file.filename.indexOf('.json') == -1) {
            fileLink.href = file.filename;
          } else {
            const parsedFeatures = JSON.stringify(file.features);
            const blob = new Blob([parsedFeatures], {type: "application/json"});
            const url  = URL.createObjectURL(blob);
            fileLink.href = url;
            // fileLink.download = file.filename;
          }
          fileLink.textContent = file.filename;
          fileLink.className = "text-gray-600";

          const listItem = document.createElement("li");
          listItem.appendChild(fileLink);

          dropdownInputChoice.appendChild(option);
          inputFilesLinks.appendChild(listItem);
          i++;
        });
        dropdownInputChoice.getElementsByTagName('option')[data.index].selected = true;
      }
    })
    .catch(error => {
      console.error("Error fetching data:", error);
      messageOutput.innerText = "An unexpected error occurred. Please try again.";
    });
}

function load_output_files() {
  fetch("/dashboard/urlsOutput")
    .then(response => response.text())
    .then(text => {
      let data = text.length > 0 ? JSON.parse(text) : {error: "Empty response"};
      if (data.error) {
        outputFilesLinks.innerHTML = data.error;
        return;
      }
      const decodedUrls = atob(data.urlsoutput);
      const urls = decodedUrls.match(/https?:\/\/[^\s]+?(?=https?:\/\/|$)/g);
      if (urls == null) {
        return;
      } else {
        outputFilesLinks.innerHTML = "";
	remove_all_children(outputFilesLinks);
        urls.forEach(url => {
          const urlLink = document.createElement("a");

          urlLink.href = url;
          urlLink.textContent = url;
          urlLink.className = "text-gray-600";

          const listItem = document.createElement("li");
          listItem.appendChild(urlLink);

          outputFilesLinks.appendChild(listItem);
        });
      }
    })
    .catch(error => {
      console.error("Error fetching data:", error);
      messageOutput.innerText = "An unexpected error occurred. Please try again.";
    });
}

window.addEventListener("load", () => {
  load_input_files();
  load_output_files();
});

page.addEventListener("click", () => {
  messageOutput.innerText = "";
});

function patch_data(form, dropdown, path) {
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    if (dropdown.disabled) {
      return;
    }
    const body = JSON.stringify({ newindex: dropdown.value })

    fetch(path, {
      method: 'PATCH',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: body
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        messageOutput.innerHTML = data.error;
      } else {
        messageOutput.innerHTML = data.message;
      }
    })
    .catch(error => {
      console.error('Error sending data:', error);
      messageOutput.innerText = 'An unexpected error occurred. Please try again.';
    });
  });
}

patch_data(formInputChoice, dropdownInputChoice, "/currentInputFile");

urls.addEventListener("submit", (e) => {
  e.preventDefault();

  const parsedUrls = textArea.value.split('\n').map(line => line.trim()).filter(line => line !== '');
  const body = JSON.stringify({ urls: parsedUrls, ai: dropdownAIChoice.value});

  messageOutput.innerText = "message sent, data is being analized (averaging 15 sec per criterion)..."

  fetch("/urls", {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: body
  })
  .then(response => response.text())
  .then(text => {
    let data = text.length > 0 ? JSON.parse(text) : {error: "Empty response"};
    if (data.error) {
      messageOutput.innerHTML = data.error;
    } else {
      handlePopup(data.message);
      messageOutput.innerHTML = "";
    }
  })
  .catch(error => {
    console.error('Error sending data:', error);
    document.getElementById("messageOutput").innerText = 'An unexpected error occurred. Please try again.';
  })
  .finally(() => {
    setTimeout(() => {
      document.getElementById("messageOutput").innerText = "";
    }, 5000);
  });
});

function showFileSelected() {
  const fileStatus = document.getElementById("fileStatus");
  const fileInput = document.getElementById("criterias");

  if (fileInput.files.length > 0) {
    fileStatus.textContent = "File selected: " + fileInput.files[0].name;
    fileStatus.classList.remove("text-blue-800");
    fileStatus.classList.add("text-blue-400");
  } else {
    fileStatus.textContent = "";
  }
}

function sendCriterias(event) {
  event.preventDefault();
  const file = document.getElementById("criterias").files[0];
  const fileStatus = document.getElementById("fileStatus");
  const submitJSON = document.getElementById("submitJSON");

  if (submitJSON.disabled) {
    return;
  }
  submitJSON.disabled = true;
  setTimeout(() => { submitJSON.disabled = false }, 3000);


  if (!file) {
    fileStatus.textContent = "Please select a file before submitting.";
    fileStatus.classList.remove("text-blue-400");
    fileStatus.classList.add("text-blue-800");
    return;
  }

  const reader = new FileReader();
  reader.onload = function() {
    const features = JSON.parse(reader.result);

    fetch('https://localhost/criterias', {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        features: features,
        filename: file.name
      }),
    })
    .then(response => response.json())
    .then(data => {
      if (!data.error) {
        fileStatus.textContent = "Criteria successfully submitted!";
        fileStatus.classList.remove("text-blue-800");
        fileStatus.classList.add("text-blue-400");
	load_input_files();
      } else {
        console.error('Error sending data:', data.error);
        fileStatus.textContent = "Error submitting criteria.";
        fileStatus.classList.remove("text-blue-400");
        fileStatus.classList.add("text-blue-800");
      }
    })
    .catch(error => {
      console.error('Error sending data:', error);
      fileStatus.textContent = "Error submitting criteria.";
      fileStatus.classList.remove("text-blue-400");
      fileStatus.classList.add("text-blue-800");
    })
    .finally(() => {
      setTimeout(() => {
        fileStatus.textContent = "";
      }, 5000);
    });
  }

  reader.readAsText(file);
}

function updateEmail(event) {
  event.preventDefault();
  const email = document.getElementById("email").value;
  const emailStatus = document.getElementById("emailStatus");

  fetch('https://localhost/updateEmail', {
    method: 'PATCH',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: email
    }),
  })
  .then(response => response.text())
  .then(data => {
    emailStatus.textContent = "Email updated successfully!";
    emailStatus.classList.remove("text-blue-800");
    emailStatus.classList.add("text-blue-400");
  })
  .catch(error => {
    console.error('Error updating email:', error);
    emailStatus.textContent = "Error updating email.";
    emailStatus.classList.remove("text-blue-400");
    emailStatus.classList.add("text-blue-800");
  })
  .finally(() => {
    setTimeout(() => {
      document.getElementById("email").value = "";
      apiKeyStatus.textContent = "";
    }, 5000);
  });
}

function updateAPIKey(event) {
  event.preventDefault();
  const apiKey = document.getElementById("apiKey").value;
  const apiKeyStatus = document.getElementById("apiKeyStatus");
  const api = apiKey.indexOf("-") == -1 ? "gemini" : "openai";
  fetch('https://localhost/updateAPIKey', {
    method: 'PATCH',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      type: api,
      apikey: apiKey
    }),
  })
  .then(response => response.text())
  .then(data => {
    apiKeyStatus.textContent = `${api} API key updated successfully!`;
    apiKeyStatus.classList.remove("text-blue-800");
    apiKeyStatus.classList.add("text-blue-400");
  })
  .catch(error => {
    console.error('Error updating API key:', error);
    apiKeyStatus.textContent = "Error updating API key.";
    apiKeyStatus.classList.remove("text-blue-400");
    apiKeyStatus.classList.add("text-blue-800");
  })
  .finally(() => {
    document.getElementById("apiKey").value = "";
    setTimeout(() => {
      apiKeyStatus.textContent = "";
    }, 5000);
  });
}

 /* exported gapiLoaded */
 /* exported gisLoaded */
 /* exported handleAuthClick */
 /* exported handleSignoutClick */


const CLIENT_ID = "726518157620-8s2194lb2ka65vfga9loee2sookpjfda.apps.googleusercontent.com";
const API_KEY = '';
const APP_ID  = "bluebeam-438322";
const SCOPES = "https://www.googleapis.com/auth/drive.file https://www.googleapis.com/auth/spreadsheets";

let tokenClient;
let accessToken = null;
let pickerInited = false;
let gisInited = false;

document.getElementById("googleSsButton").style.visibility = "hidden";

function gapiLoaded() {
  gapi.load('client:picker', initializePicker);
}

/**
 * Callback after the API client is loaded. Loads the
 * discovery doc to initialize the API.
 */
async function initializePicker() {
  await gapi.client.load('https://www.googleapis.com/discovery/v1/apis/drive/v3/rest');
  await gapi.client.load('https://sheets.googleapis.com/$discovery/rest?version=v4');
  pickerInited = true;
  maybeEnableButtons();
}

/**
 * Callback after Google Identity Services are loaded.
 */
function gisLoaded() {
  tokenClient = google.accounts.oauth2.initTokenClient({
    client_id: CLIENT_ID,
    scope: SCOPES,
    callback: '', // defined later
  });
  gisInited = true;
  maybeEnableButtons();
}

function maybeEnableButtons() {
  if (pickerInited && gisInited) {
    document.getElementById('googleSsButton').style.visibility = 'visible';
  }
}


function initOAuth(event) {
   tokenClient.callback = async (response) => {
      if (response.error !== undefined) {
        throw (response);
      }
      accessToken = response.access_token;
      await createPicker();
   };

   if (accessToken === null) {
    // Prompt the user to select a Google Account and ask for consent to share their data
    // when establishing a new session.
     tokenClient.requestAccessToken({prompt: 'consent'});
   } else {
    // Skip display of account chooser and consent dialog for an existing session.
     tokenClient.requestAccessToken({prompt: ''});
   }
}

function createPicker() {
  const view = new google.picker.View(google.picker.ViewId.SPREADSHEETS);
  const picker = new google.picker.PickerBuilder()
      .enableFeature(google.picker.Feature.NAV_HIDDEN)
      // .enableFeature(google.picker.Feature.MULTISELECT_ENABLED)
      .setDeveloperKey(API_KEY)
      .setAppId(APP_ID)
      .setOAuthToken(accessToken) // Use the token from the server
      .addView(view)
      .addView(new google.picker.DocsUploadView())
      .setCallback(pickerCallback)
      .build();
  picker.setVisible(true);
}

async function pickerCallback(data) {
  if (data[google.picker.Response.ACTION] === google.picker.Action.PICKED) {
    const document = data[google.picker.Response.DOCUMENTS][0];
    const fileId = document[google.picker.Document.ID];

    const res = await gapi.client.sheets.spreadsheets.values.get({
      spreadsheetId: fileId,
      range: 'Sheet1', // Adjust the sheet name and range as needed
    });

    const rows = res.result.values;

    if (!rows || rows.length === 0) {
      console.log('No data found in the spreadsheet.');
      return;
    }

    // Convert rows into JSON objects assuming the first row contains headers
    const headers = rows[0];
    const jsonArray = rows.slice(1).map(row => {
      let obj = {};
      headers.forEach((header, index) => {
        obj[header] = row[index] || "";
      });
      return obj;
    });

    fetch('https://localhost/criterias', {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        features: jsonArray,
        filename: ("https://docs.google.com/spreadsheets/d/" + fileId)
      }),
    })
    .then(response => response.json())
    .then(data => {
      if (!data.error) {
        fileStatus.textContent = "Criteria successfully submitted!";
        fileStatus.classList.remove("text-blue-800");
        fileStatus.classList.add("text-blue-400");
	load_input_files();
      } else {
        console.error('Error sending data:', data.error);
        fileStatus.textContent = "Error submitting criteria.";
        fileStatus.classList.remove("text-blue-400");
        fileStatus.classList.add("text-blue-800");
      }
    })
    .catch(error => {
      console.error('Error sending data:', error);
      fileStatus.textContent = "Error submitting criteria.";
      fileStatus.classList.remove("text-blue-400");
      fileStatus.classList.add("text-blue-800");
    });
  }
}
