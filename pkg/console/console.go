package console

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/ariary/go-utils/pkg/clipboard"
	"github.com/ariary/go-utils/pkg/color"
	encryption "github.com/ariary/go-utils/pkg/encrypt"
	"github.com/ariary/quicli/pkg/quicli"
)

var cmdDir = ""

func Console(flags quicli.Config) {
	// flag & var
	var embed bool
	//detect if we are in remote mode
	if EmbedCert != "" && EmbedKey != "" {
		fmt.Println(color.Italic(color.Dim("Key and Cert embed in binary detected")))
		embed = true
	}

	port := ":" + flags.GetStringFlag("port")
	addr := flags.GetStringFlag("url")
	url := addr + port

	var wsEndpoint string
	if flags.GetBoolFlag("secure") {
		wsEndpoint = encryption.GenerateRandom()
	} else {
		wsEndpoint = "sh"
	}
	//load current directory
	cmdDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// launch server
	fmt.Println("🚀 Launch 'console.sh' websocket server listening on", color.Italic(color.Yellow(port)))

	//log info
	fmt.Println("📁 Serve on directory:", color.Italic(color.Yellow(cmdDir)))
	fmt.Println()
	fmt.Println("📋 Copy paste in browser console:")
	command := "s=new WebSocket(\"wss://" + url + "/" + wsEndpoint + "\"),s.onmessage=function(ev){console.log(ev.data)};function sh(cmd){s.send(cmd)};function promptsh(){cmd=prompt();s.send(cmd)};Object.defineProperty(window, 'psh', { get: promptsh });"
	fmt.Println(color.Teal(command))
	clipboard.Copy(command)
	fmt.Println("👀 Or simply visit:")
	fmt.Println(color.Teal("https://" + url))
	fmt.Println("💻 For a \"in-web-terminal\"")
	fmt.Println(color.Teal("https://" + url + "/interactive"))
	fmt.Println()

	setupRoutes(command, url, wsEndpoint)

	// launch webserver
	if embed {
		//write file
		if err := os.WriteFile("cert.pem", []byte(EmbedCert), 0400); err != nil {
			fmt.Println("Failed writing embedding cert.pem:")
		}
		if err := os.WriteFile("key.pem", []byte(EmbedKey), 0400); err != nil {
			fmt.Println("Failed writing embedding key.pem:")
		}
	}
	err = http.ListenAndServeTLS(port, "cert.pem", "key.pem", nil)
	if err != nil {
		// try to generate cert
		if strings.Contains(err.Error(), "no such file or directory") {
			//try to generate cert
			if errMkcert := generateCert(addr, flags.GetBoolFlag("privileged")); errMkcert != nil {
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
