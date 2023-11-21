package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB

//trying to recover todays hardwork

func init() {

	// Open a database connection
	var err error
	db, err = sql.Open("mysql", "root:Xonen@3616@tcp(127.0.0.1:3306)/twitter_clone")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	log.Println("Connected to the database")
}

func handleCreateTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusInternalServerError)
			return
		}

		// Get tweet content from the form
		tweetContent := r.FormValue("tweetContent")

		// Save the tweet content to the database
		_, err = db.Exec("INSERT INTO tweets (content) VALUES (?)", tweetContent)
		if err != nil {
			log.Println("Error storing tweet in the database:", err)
			http.Error(w, "Error storing tweet in the database", http.StatusInternalServerError)
			return
		}

		// Respond to the client with a JSON response
		response := map[string]string{
			"status":  "success",
			"message": fmt.Sprintf("Tweet created: %s", tweetContent),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create_tweet", handleCreateTweet).Methods("POST")

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Update this with your frontend's actual origin
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Use the CORS middleware
	handler := c.Handler(r)

	http.Handle("/", handler)

	http.ListenAndServe(":8080", nil)
}
