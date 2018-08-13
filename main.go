// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseFiles("root.html"))

// Server holds database and other information about this server.
type Server struct {
	db *sql.DB
}

func (x *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "root.html", "nodata")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (x *Server) viewHandler(w http.ResponseWriter, r *http.Request) {
	pass := r.PostFormValue("pass")
	hash := fmt.Sprintf("%X", sha1.Sum([]byte(pass)))

	rows, err := x.db.Query("SELECT count FROM pwned where hash = \"" + hash + "\"")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	count := 0
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "{ \"count\":%d }\n", count)
}

func main() {
	var err error

	srv := &Server{}

	srv.db, err = sql.Open("sqlite3", "/data/tmp/pwned/pwned.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer srv.db.Close()

	http.HandleFunc("/", srv.rootHandler)
	http.HandleFunc("/view/", srv.viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
