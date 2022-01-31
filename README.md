# console.sh

Execute shell command from browser's console (developer tools)

## Why?

Why not! The need does not inspire the feature, it's the other way around *(s/o Apple philosophy)*

## Use

Launch websocket server with certificates (otherwise browsers won't accept connection):
```shell
mkcert -install
mkcert -key-file key.pem -cert-file cert.pem localhost 127.0.0.1 ::1 # Many way to do it, openssl etc => key: key.pem and cert: cert.pem
./console.sh-server
```

Open browser's console (certainly with `Shift + CTRL + J` or `Shift + ⌘ + J`). Copy/paste within:
```
s=new WebSocket("wss://localhost:8080/sh"),s.onmessage=function(ev){console.log(ev.data)}
```

Now execute shell command from browser console with: `s.send([cmd])`

## Notes
* SOP and CORS don't apply to websocket, **However** CSP does. Many websites specify `connect-src` CSP directive which prevent restrict loaded URL from WebSocket (⇒ can't use `conbsole.sh` on these websites)
* Without `wss` (secure websocket) browser wuldn't authorize websocket communication ⇒ need certificate and key
* **⚠️ This project is not secure! Use it with parsimony and of course shoot out the server when when you are done using it**
