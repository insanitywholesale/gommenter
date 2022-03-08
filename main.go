package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const createtable = `CREATE TABLE IF NOT EXISTS comments (name TEXT, content TEXT);`

var (
	pgClient   *sql.DB
	thankPage  string
	commitHash string
	commitDate string
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
		log.Println(name, "said:", content)
		// redirect to thank you page
		log.Println("thankpage", thankPage)
		http.Redirect(w, r, thankPage, 301)
	}
	return
}

func getComments(w http.ResponseWriter, r *http.Request) {
	var results string
	var curRowName string
	var curRowContent string
	if r.Method == "GET" {
		privID := r.Header.Get("Comrade-ID")
		log.Println("privID:", privID)
		comID := os.Getenv("COMRADE_ID")
		if privID == comID {
			log.Println("comradeID:", comID)
			query := `SELECT name, content FROM comments`
			rows, err := pgClient.Query(query)
			if err != nil {
				log.Println("uh-oh, getComments stinky!")
				return
			}
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&curRowName, &curRowContent)
				if err != nil {
					log.Println(err)
				}
				results = results + "user: " + curRowName + " said: " + curRowContent + "\n"
			}
			log.Print("comments-> ", results)
			w.Write([]byte(results))
			return
		}
		return
	}
	return
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte("commitHash: " + commitHash + "\n"))
		w.Write([]byte("commitDate: " + commitDate + "\n"))
		return
	}
	return
}

func main() {
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		pgURL = "postgresql://tester:Apasswd@localhost:5432/?sslmode=disable"
	}
	client, err := newPostgresClient(pgURL)
	if err != nil {
		log.Fatal("error connecting to postgres", err)
	}
	pgClient = client
	thankPage = os.Getenv("THANK_PAGE")
	if thankPage == "" {
		thankPage = "http://localhost:1313/thanks"
	}
	http.HandleFunc("/api/form/new", addComment)
	http.HandleFunc("/api/comments", getComments)
	http.HandleFunc("/info", getInfo)
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
