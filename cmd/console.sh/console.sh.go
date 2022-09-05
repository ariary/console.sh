package main

import (
	"github.com/ariary/console.sh/pkg/console"
	"github.com/ariary/quicli/pkg/quicli"
)

func main() {

	cli := quicli.Cli{
		Usage:       "console.sh [flags]",
		Description: "Share your terminal in your browser console",
		Flags: quicli.Flags{
			{Name: "url", Default: "localhost", Description: "Websocket server URL"},
			{Name: "port", Default: "8080", Description: "Websocket server port"},
			{Name: "secure", Default: false, Description: "Protect websocket endpoint with random name"},
			{Name: "privileged", Default: true, Description: "Define if user is privileged. Used to install CA in the system trust store"},
		},
		Function: console.Console,
	}
	cli.Run()
}
