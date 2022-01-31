# console.sh

Execute shell command from browser console (developer tools)

## Usage

Launch websocket server with certificates (otherwise browsers won't accept connection):
```shell
mkcert -install
mkcert -key-file key.pem -cert-file cert.pem localhost 127.0.0.1 ::1 # Many way to do it, openssl etc => key: key.pem and cert: cert.pem
./console.sh-server
```

Open browser console (certainly with `Shift + CTRL + J` or `Shift + ⌘ + J`). Copy/paste within:
```javascript
s=new WebSocket("wss://localhost:8080/sh"),s.onmessage=function(ev){console.log(ev.data)};function sh(cmd){s.send(cmd)};
```

Now execute shell command from browser console with: `sh("[cmd]")`

## Why?

Why not! The need does not inspire the feature, it's the other way around *(s/o Apple philosophy)*

## Notes
* SOP and CORS don't apply to websocket, **However** CSP does. Many websites specify `connect-src` CSP directive which prevent restrict loaded URL from WebSocket (⇒ can't use `conbsole.sh` on these websites, empty new tabs will do the job)
* Without `wss` (secure websocket) browser wouldn't authorize websocket communication ⇒ need certificate and key
* **⚠️ This project is not secure! Use it with parsimony and of course shoot out the server when when you are done using it**
