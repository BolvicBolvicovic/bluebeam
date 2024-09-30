document.getElementById('scrapeButton').addEventListener('click', function() {
  const username = document.getElementById('username').value;
  const sessionKey = document.getElementById('sessionKey').value;

  if (!username || !sessionKey) {
    alert('Please provide both username and session key.');
    return;
  }

  // Get the active tab and run the scraping script with username and sessionKey
  browser.tabs.query({ active: true, currentWindow: true }).then((tabs) => {
    browser.tabs.executeScript(tabs[0].id, {
      code: '(' + scrapeAndSendContent.toString() + ')("' + username + '", "' + sessionKey + '");'
    });
  });
});

// Scraping function that runs in the context of the webpage
function scrapeAndSendContent(username, sessionKey) {
  // Scrape the content
  let links = Array.from(document.querySelectorAll('a')).map(a => a.href);
  let buttons = Array.from(document.querySelectorAll('button')).map(b => ({
    text: b.innerText,
    onclick: b.onclick ? b.onclick.toString() : null
  }));
  let pageHtml = document.documentElement.outerHTML;

  // Send the scraped data along with the username and session key to your Go server
  fetch('http://localhost:8080/analyze', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, sessionKey, links, buttons, pageHtml })
  })
  .then(response => response.json())
  .then(data => {
    console.log('Analysis result:', data);
  })
  .catch(error => console.error('Error sending data:', error));
}
