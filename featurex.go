package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"github.com/alessio/shellescape"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func MyHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := session.Values["email"]

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

	// ok: gorilla-command-injection-taint
	cmd3 := exec.Command("bash")
	cmd3Writer, _ := cmd3.StdinPipe()
	cmd3.Start()
	cmd3Input := fmt.Sprintf("ls %s", email)

	// ruleid: gorilla-command-injection-taint
	cmd3Writer.Write([]byte(cmd3Input + "\n"))

	// ruleid: gorilla-command-injection-taint
	io.WriteString(cmd3Writer, cmd3Input)

	cmd4Input := shellescape.Quote(email)
	// ok: gorilla-command-injection-taint
	syscall.Exec("echo " + cmd4Input)
}

//Gorilla schema

var decoder = schema.NewDecoder()

type Person struct {
	Name  string
	Phone string
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := decoder.Decode(&person, r.PostForm)

	cmd := &exec.Cmd{
		// Path is the path of the command to run.
		// ruleid: gorilla-command-injection-taint
		Path: person.Name,
		// Args holds command line arguments, including the command as Args[0].
		Args:   []string{"tr", "--help"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd.Start()
	cmd.Wait()
}

// Gorilla securecookie

func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {

	// Hash keys should be at least 32 bytes long
	var hashKey = []byte("very-secret")
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte("a-lot-secret")
	var s = securecookie.New(hashKey, blockKey)

	if cookie, err := r.Cookie("cookie-name"); err == nil {

		value := make(map[string]string)

		if err = s.Decode("cookie-name", cookie.Value, &value); err == nil {

			cmd := &exec.Cmd{
				// Path is the path of the command to run.
				// ruleid: gorilla-command-injection-taint
				Path: value["foo"],
				// Args holds command line arguments, including the command as Args[0].
				Args:   []string{"tr", "--help"},
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			cmd.Start()
			cmd.Wait()
		}

		var cookies = map[string]*securecookie.SecureCookie{
			"previous": securecookie.New(
				securecookie.GenerateRandomKey(64),
				securecookie.GenerateRandomKey(32),
			),
			"current": securecookie.New(
				securecookie.GenerateRandomKey(64),
				securecookie.GenerateRandomKey(32),
			),
		}
		err = securecookie.DecodeMulti("cookie-name", cookie.Value, &value, cookies["current"], cookies["previous"])

		cmd := &exec.Cmd{
			// Path is the path of the command to run.
			// ruleid: gorilla-command-injection-taint
			Path: value["foo"],
			// Args holds command line arguments, including the command as Args[0].
			Args:   []string{"tr", "--help"},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		cmd.Start()
		cmd.Wait()

	}
}

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
