package main

import (
	"net/http"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// gorilla sessions

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func MyHandler2(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := session.Values["email"]
	opt, err := pg.ParseURL("postgres://user:pass@localhost:5432/db_name")
	db := pg.Connect(opt)

	// ruleid: gorilla-pg-sqli-taint
	res, err := db.Prepare("SELECT name FROM users WHERE email=?" + email)

}

// gorilla/schema
var decoder = schema.NewDecoder()

type Person struct {
	Name  string
	Phone string
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := decoder.Decode(&person, r.PostForm)

	opt, err := pg.ParseURL("postgres://user:pass@localhost:5432/db_name")
	db := pg.Connect(opt)

	// ruleid: gorilla-pg-sqli-taint
	res, err := db.Prepare("SELECT name FROM users WHERE email=?" + person.Name)

	return nil
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
		err = s.Decode("cookie-name", cookie.Value, &value)

		opt, err := pg.ParseURL("postgres://user:pass@localhost:5432/db_name")
		db := pg.Connect(opt)

		// ruleid: gorilla-pg-sqli-taint
		res, err := db.Prepare("SELECT name FROM users WHERE email=?" + value["email"])

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

		// ruleid: gorilla-pg-sqli-taint
		res, err := db.Prepare("SELECT name FROM users WHERE email=?" + value["email"])

	}
}
