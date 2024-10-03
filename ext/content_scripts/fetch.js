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
	// Scrape the content
  	let links = Array.from(document.querySelectorAll('a')).map(a => a.href);
  	let buttons = Array.from(document.querySelectorAll('button')).map(b => ({
  	  text: b.innerText,
  	  onclick: b.onclick ? b.onclick.toString() : null
  	}));
  	let pageHtml = document.documentElement.outerHTML;

  	// Send the scraped data along with the username and session key to your Go server
  	fetch('https://localhost:443/analyze', {
  	  method: 'POST',
  	  headers: { 'Content-Type': 'application/json' },
  	  body: JSON.stringify({ 
	    username: message.username,
	    sessionkey: message.sessionKey,
	    links,
	    buttons,
	    pageHtml
	  })
  	})
  	.then(response => response.json())
  	.then(data => {
  	  console.log('Analysis result:', data);
  	})
  	.catch(error => console.error('Error sending data:', error));
  });
})();
