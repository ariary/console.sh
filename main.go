package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/ariary/go-utils/pkg/check"
	"github.com/ariary/go-utils/pkg/clipboard"
	"github.com/ariary/go-utils/pkg/color"
	"github.com/gorilla/websocket"
)

const homePageTpl = `
<!DOCTYPE html>
<html>
<title>console.sh: Home</title>
<body>
	<h1>console.sh: Home Page</h1>
	<p>Open browser console and interact with console.sh with:<br>
	<pre><code>
	> sh(\"[command]\")
	//OR (prompted version)
	> psh
	</code></pre>
	<script>{{ .ConnectionScript }}</script>
	<br><br>
	<p> In other tab, first connect to the console.sh websocket server with:<br>
	<pre><code>
	{{ .ConnectionScript}}
	</code></pre>
</body>
</html>
`

const interactivePageTpl = `
<html>
<title>interactive console.sh</title>
<body>
	<h1>Enter command in box</h1>
	<label for="command">Enter your command:</label>
	<script>
	//result listener
	s=new WebSocket("wss://{{ .Url}}/sh"),s.onmessage=function(ev){document.getElementById("result").innerHTML=ev.data};
	function sendCommand(){
		console.log("toto")
		cmd = document.getElementById("command").value
		console.log(cmd)
		s.send(cmd)
	}
	</script>
	<form action="javascript:;" onsubmit="sendCommand(this)">
	<input name='command' type='text' id='command'>
        <button id='btn' class='btn btn-primary' type='submit'>execute</button>
	</form>
	<div id="result">
	</div>
</body>
</html>
`

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //do not use it if you want to construct a robust websocket server
}

var cmdDir = ""

//execute command on shell and return stdout & stderr
func execute(cmd string) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := exec.Command("/bin/sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr
	command.Dir = cmdDir
	err := command.Run()
	if err != nil {
		fmt.Println("failed execute command", cmd, ":", err)
	}
	if stderr.String() != "" {
		return stdout.String() + stderr.String()
	}

	return stdout.String()

}

// reader: read message from websocket
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		cmd := string(p)
		log.Println("command from client:", cmd)

		//Manage special cases
		if strings.HasPrefix(cmd, "cd") { //change directory where execute command
			args := strings.Split(cmd, " ")
			if len(args) < 2 {
				cdErrorMsg := []byte("console.sh: 'cd' need an argument")
				if err := conn.WriteMessage(messageType, cdErrorMsg); err != nil {
					log.Println(err)
					return
				}
			}
			cmdDir = args[1]
			fmt.Println("client change command directory:", cmdDir)
		} else if strings.HasPrefix(cmd, "exit") { //exit shell
			exitMsg := []byte("console.sh: exit console.sh (close websocket)")
			if err := conn.WriteMessage(messageType, exitMsg); err != nil {
				log.Println(err)
				return
			}
			conn.Close()
		} else {
			//Execute command and return result
			output := execute(cmd)
			responseBytes := []byte(output)
			if err := conn.WriteMessage(messageType, responseBytes); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

//homePage: home handler
func homePage(script string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remote := strings.Split(r.RemoteAddr, ":")[0]
		fmt.Println(color.Purple(remote), "Visit home page")
		t, err := template.New("homepage").Parse(homePageTpl)
		check.Check(err, "failed loading home template")
		data := struct {
			ConnectionScript string
		}{
			ConnectionScript: script,
		}

		check.Check(t.Execute(w, data), "failed writing script in home page")
	})
}

//interactivePage: interactive handler
func interactivePage(url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remote := strings.Split(r.RemoteAddr, ":")[0]
		fmt.Println(color.Purple(remote), "Visit interactive terminal page")
		t, err := template.New("interactive").Parse(interactivePageTpl)
		check.Check(err, "failed loading interactive template")
		data := struct {
			Url string
		}{
			Url: url,
		}

		check.Check(t.Execute(w, data), "failed writing script in interactive page")
	})
}

//wsEndpoint: Handler for /sh endpoint. Websocket connection
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Client connected to console.sh!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws) // listen indefinitely for new messages coming through on our bebSocket connection
}

func setupRoutes(script string, url string) {
	//http.HandleFunc("/", homePage)
	http.Handle("/", homePage(script))
	http.Handle("/interactive", interactivePage(url))
	http.HandleFunc("/sh", wsEndpoint)
}

//generateCert: try to generate cert in current directory with mkcert.
func generateCert(addr string) error {
	// Check for mkcert installation
	if _, err := exec.LookPath("mkcert"); err != nil {
		return err
	}

	//create local CA and install it
	mkcertInstallArgs := []string{"mkcert", "-install"}
	mkcertInstallCmd := exec.Command(mkcertInstallArgs[0], mkcertInstallArgs[1:]...)

	if err := mkcertInstallCmd.Start(); err != nil {
		return err
	}

	if err := mkcertInstallCmd.Wait(); err != nil {
		return err
	}

	mkcertGenerateArgs := []string{"mkcert", "--key-file", "key.pem", "-cert-file", "cert.pem", addr, "127.0.0.1", "::1"}
	mkcertGenerateCmd := exec.Command(mkcertGenerateArgs[0], mkcertGenerateArgs[1:]...)

	if err := mkcertGenerateCmd.Start(); err != nil {
		return err
	}

	return mkcertGenerateCmd.Wait()
}

func main() {
	// flag & var
	var addr, p string
	flag.StringVar(&addr, "url", "localhost", "Websocket server URL")
	flag.StringVar(&p, "p", "8080", "Websocket server port")
	flag.Parse()
	port := ":" + p

	//launch server
	fmt.Println("🚀 Launch 'console.sh' websocket server listening on", color.Italic(color.Yellow(port)))
	//Load current directory
	cmdDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	//log info
	url := addr + port
	fmt.Println("📁 Serve on directory:", color.Italic(color.Yellow(cmdDir)))
	fmt.Println()
	fmt.Println("📋 Copy paste in browser console:")
	command := "s=new WebSocket(\"wss://" + url + "/sh\"),s.onmessage=function(ev){console.log(ev.data)};function sh(cmd){s.send(cmd)};function promptsh(){cmd=prompt();s.send(cmd)};Object.defineProperty(window, 'psh', { get: promptsh });"
	fmt.Println(color.Teal(command))
	clipboard.Copy(command)
	fmt.Println("👀 Or simply visit:")
	fmt.Println(color.Teal("https://" + url))
	fmt.Println("💻 For a \"in-web-terminal\"")
	fmt.Println(color.Teal("https://" + url + "/interactive"))
	fmt.Println()

	setupRoutes(command, url)

	// launch webserver
	err = http.ListenAndServeTLS(port, "cert.pem", "key.pem", nil)
	if err != nil {
		// try to generate cert
		if strings.Contains(err.Error(), "no such file or directory") {
			//try to generate cert
			if errMkcert := generateCert(addr); errMkcert != nil {
				fmt.Println(color.Evil("Failed to generate cert with mkcert", errMkcert))
				os.Exit(1)
			}
			fmt.Println("ℹ️ Generate cert with mkcert in current directory (cert.pem and key.pem).. Restart server")
			log.Fatal(http.ListenAndServeTLS(port, "cert.pem", "key.pem", nil))
		} else {
			fmt.Println(color.Evil("Failed to start server:"), err)
			os.Exit(1)
		}
	}
}
