function listenForClicks() {
  document.addEventListener("click", (e) => {
    function reportError(error) {
      console.error(`Error caught: ${error}`);
    }
    const username = document.getElementById('username').value;
    const sessionKey = document.getElementById('sessionKey').value;
    console.log("Scrape button clicked")

    if (!username || !sessionKey) {
      alert('Please provide both username and session key.');
      return;
    }
    browser.tabs
      .query({ active: true, currentWindow: true })
      .then((tabs) => {
        browser.tabs.sendMessage( tabs[0].id, {
          username: username,
          sessionKey: sessionKey,
        });
      })
      .catch(reportError);
  });
}

browser.tabs
  .executeScript({ file: "/content_scripts/fetch.js" })
  .then(listenForClicks)
