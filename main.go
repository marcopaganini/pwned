// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha1"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	rootTemplate = template.Must(template.ParseFiles("root.html"))

	// SHA1 matching regexp
	sha1Regex = regexp.MustCompile(`(?i)[\da-f]{40}`)
)

// Server holds database and other information about this server.
type Server struct {
	db *sql.DB
}

func (x *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	err := rootTemplate.ExecuteTemplate(w, "root.html", "nodata")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (x *Server) viewHandler(w http.ResponseWriter, r *http.Request) {
	var hash string

	// Fetch password from POST request and calculate the uppercase
	// (textual) version of the SHA1 hash. If the password looks like
	// a hash (40 hexascii chars), use it directly.
	pass := r.PostFormValue("pass")
	if sha1Regex.MatchString(pass) {
		hash = strings.ToUpper(pass)
	} else {
		hash = fmt.Sprintf("%X", sha1.Sum([]byte(pass)))
	}

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

	dbfile := flag.String("dbfile", "", "SQLite3 Database file.")
	rootpath := flag.String("rootpath", "", "Root path in the URL (usually empty).")
	port := flag.Int("port", 8080, "HTTP server port.")

	flag.Parse()

	srv := &Server{}

	srv.db, err = sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer srv.db.Close()

	http.HandleFunc(*rootpath+"/", srv.rootHandler)
	http.HandleFunc(*rootpath+"/view/", srv.viewHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
