package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const createtable = `CREATE TABLE IF NOT EXISTS comments (name TEXT, content TEXT);`

func addComment(w http.ResponseWriter, r *http.Request) {
	// this db stuff is here because otherwise there is a nil pointer dereference
	// unfortunate but such is life
	postgresURL := os.Getenv("PG_URL")
	if postgresURL == "" {
		postgresURL = "postgres://tester:Apasswd@localhost:5432?sslmode=disable"
	}
	pgClient, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Println("error open", err)
	}
	err = pgClient.Ping()
	if err != nil {
		log.Fatal("error ping", err)
	}
	// create database if it doesn't exist
	_, err = pgClient.Exec(createtable)
	if err != nil {
		log.Println("error exec", err)
	}

	// actually handle request
	// but only if method is post
	if r.Method == "POST" {
		var name string
		var content string
		err := r.ParseForm()
		if err != nil {
			log.Fatal("error parseform", err)
			return
		}
		if (r.Form["name"] != nil) && (len(r.Form["name"][0]) > 0) {
			name = r.Form["name"][0]
		} else {
			name = "anon"
		}
		if r.Form["content"] != nil {
			content = r.Form["content"][0]
		} else {
			return
		}
		// save comment to database
		query := `INSERT INTO comments (name, content) VALUES ($1, $2)`
		_, err = pgClient.Exec(query, name, content)
		if err != nil {
			return
		}
		// make a redirect to r.URL.Host contact thank you page
		thankPage := os.Getenv("THANK_PAGE")
		if thankPage == "" {
			thankPage = "http://localhost:1313/thanks"
		}
		http.Redirect(w, r, thankPage, 301)
		return
	}
	return
}

func main() {
	http.HandleFunc("/api/form/new", addComment)
	err := http.ListenAndServe(":9097", nil)
	if err != nil {
		log.Fatal("error starting http server", err)
		return
	}
}
