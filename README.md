# Backtone
Backtone is a program that converts websites to Atom RSS
feeds and prints the results to stdout. It should be called via a
server-side script that then writes the resulting XML to the webpage.

The name is a playful nod toward converting contents back to the web of yesteryear.
Why suffer in Web 2/3 when you could go "back to one"?

# Dependencies
To ensure Backtone always supports current ECMAScript standards, it needs a headless browser instance to execute client-side scripts. To avoid Cloudflare/similar challenge interference, we use [FlareSolverr](https://developer.chrome.com/blog/chrome-headless-shell). Its Docker image also contains the browser we use.

# Usage
All files needed for usage, including an the input XML, server-side script, and reverse proxy config, have examples in the `examples` folder in this repo. 

1. Create a server-side script that calls `backtone` with the only argument being a base64-encoded XML in the format specified when running `backbone` with no arguments
2. Serve that script via a reverse proxy
    - Now when a feed reader fetches the feed, the feed will be created in real time
