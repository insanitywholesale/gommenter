package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Comment struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

const createtable = `CREATE TABLE IF NOT EXISTS comments (name TEXT, content TEXT);`

var pgClient *sql.DB

func newPostgresClient(postgresURL string) (*sql.DB, error) {
	// create postgres client
	client, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	err = client.Ping()
	if err != nil {
		return nil, err
	}
	_, err = client.Exec(createtable)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func addComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var formErrors = []string{}
		var name string
		var content string
		r.ParseForm()
		if (r.Form["name"] != nil) && (len(r.Form["name"][0]) > 0) {
			name = r.Form["name"][0]
		} else {
			formErrors = append(formErrors, "imma need ya name, sweaty")
		}
		if r.Form["content"] != nil {
			content = r.Form["content"][0]
		} else {
			formErrors = append(formErrors, "so you have nothing to say?")
		}
		// TODO: make this error handling better
		if len(formErrors) > 0 {
			return
		}
		// save comment to database if there are no errors
		query := `INSERT INTO comments (name, content) VALUES ($1, $2);`
		_, err := pgClient.Exec(query, name, content)
		// TODO: make this error handling better
		if err != nil {
			return
		}
		// make a redirect to r.URL.Host contact thank you page
		// http.Redirect(w, r, r.URL.Host, 301)
	}
	return
}

func main() {
	pgURL := os.Getenv("POSTGRES_URL")
	pgClient, err := newPostgresClient(pgURL)
	if err != nil {
		log.Println("pgClient", pgClient)
		log.Fatal("error connecting to postgres", err)
	}
	http.HandleFunc("/api/form/new", addComment)
	http.ListenAndServe(":9097", nil)
}
