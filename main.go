package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const createtable = `CREATE TABLE IF NOT EXISTS comments (name TEXT, content TEXT);`

var (
	pgClient *sql.DB
	thankPage string
)

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
		// redirect to thank you page
		log.Println("thankpage", thankPage)
		http.Redirect(w, r, thankPage, 301)
	}
	return
}

func main() {
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		pgURL = "postgresql://tester:Apasswd@localhost:5432/?sslmode=disable"
	}
	client, err := newPostgresClient(pgURL)
	pgClient = client
	if err != nil {
		log.Fatal("error connecting to postgres", err)
	}
	thankPage = os.Getenv("THANK_PAGE")
	if thankPage == "" {
		thankPage = "http://localhost:1313/thanks"
	}
	http.HandleFunc("/api/form/new", addComment)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9097"
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("error starting http server", err)
		return
	}
}
