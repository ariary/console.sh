# console.sh

<div align=center>
<img src=https://github.com/ariary/console.sh/blob/main/img/logo.png width=150>


  <b>Execute shell command from browser console (developer tools)</b>
</div>


## Usage

Launch websocket server:
```shell
console.sh
```

Open browser console (certainly with `Shift + CTRL + K` or `Shift + ⌘ + K`). Copy/paste within:
```javascript
s=new WebSocket("wss://localhost:8080/sh"),s.onmessage=function(ev){console.log(ev.data)};function sh(cmd){s.send(cmd)};function promptsh(){cmd=prompt();s.send(cmd)};Object.defineProperty(window, 'psh', { get: promptsh });
```


Now you are able to execute shell command from browser console with:
```javascript
> sh("[command]")
//OR (alternative)
> sh`[command]`
//OR (prompted version)
> psh
```
<div align=center><img src=https://github.com/ariary/console.sh/blob/main/console.sh.png></div>

**N.B:** 
* It is also possible to **connect to a remote shell** (`localhost` is not enforced). To do so use `--url` flag. 
* **Visit home page is like a light [`gotty`](https://github.com/yudai/gotty)**. The home page automatically load the shell in your browser console


## Why?

Why not! The need does not inspire the feature, it's the other way around *(s/o Apple philosophy)*

## Set-up

Install `console.sh`:
```shell
curl -lO -L https://github.com/ariary/console.sh/releases/latest/download/console.sh
chmod +x console.sh
# or with go
go install github.com/ariary/console.sh@latest
```

Then as you have to launch the websocket server with certificates (otherwise browsers won't accept connection). Create cert and key in the same directory:
```shell
mkcert -install
mkcert -key-file key.pem -cert-file cert.pem localhost 127.0.0.1 ::1 # Many way to do it, openssl etc => key: key.pem and cert: cert.pem
```

## Notes
* SOP and CORS don't apply to websocket, **However** CSP does. Many websites specify `connect-src` CSP directive which restricts loaded URL from WebSocket (⇒ can't use `console.sh` on these websites, empty new tabs will do the job)
* Without `wss` (secure websocket) browser wouldn't authorize websocket communication ⇒ need certificate and key
* **⚠️ This project is not secure! Use it with parsimony and of course shut down the server when you are done using it**. Otherwise an XSS can easily become a RCE
