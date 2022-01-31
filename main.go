package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

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
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "console.sh: Home Page")
	//TO DO: html page with script that initiate websocket => light gotty
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

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/sh", wsEndpoint)
}

func main() {
	port := ":8080"
	fmt.Println("Launch 'console.sh' websocket server listening on", port)
	setupRoutes()
	log.Fatal(http.ListenAndServeTLS(port, "cert.pem", "key.pem", nil))

	//Load current directory
	cmdDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Serve on directory:", cmdDir)
	fmt.Println("Copy paste in browser console:")
	fmt.Println("s=new WebSocket(\"wss://localhost:8080/sh\"),s.onmessage=function(ev){console.log(ev.data)}")
}
