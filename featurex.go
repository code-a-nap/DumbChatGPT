package main

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// gorilla websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	messageType, p, err := conn.ReadMessage()
	email := string(p)

	cmd := &exec.Cmd{
		// Path is the path of the command to run.
		// ruleid: gorilla-command-injection-taint
		Path: email,
		// Args holds command line arguments, including the command as Args[0].
		Args:   []string{"tr", "--help"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd.Start()
	cmd.Wait()

	messageType, reader, err := conn.NextReader()
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	email = string(buf)

	cmd = &exec.Cmd{
		// Path is the path of the command to run.
		// ruleid: gorilla-command-injection-taint
		Path: email,
		// Args holds command line arguments, including the command as Args[0].
		Args:   []string{"tr", "--help"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd.Start()
	cmd.Wait()
}
