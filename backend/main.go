package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/rs/cors"
)

var db *sql.DB

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

var tweetType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tweet",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"tweets": &graphql.Field{
				Type: graphql.NewList(tweetType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					rows, err := db.Query("SELECT id, content FROM tweets")
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var tweets []map[string]interface{}
					for rows.Next() {
						var id int
						var content string
						err := rows.Scan(&id, &content)
						if err != nil {
							return nil, err
						}
						tweet := map[string]interface{}{
							"id":      id,
							"content": content,
						}
						tweets = append(tweets, tweet)
					}

					return tweets, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createTweet": &graphql.Field{
				Type:        tweetType,
				Description: "Create a new tweet",
				Args: graphql.FieldConfigArgument{
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					content, _ := params.Args["content"].(string)
					result, err := db.Exec("INSERT INTO tweets (content) VALUES (?)", content)
					if err != nil {
						return nil, err
					}

					id, _ := result.LastInsertId()

					createdTweet := map[string]interface{}{
						"id":      id,
						"content": content,
					}

					return createdTweet, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	var result *graphql.Result
	query := r.URL.Query().Get("query")
	if query != "" {
		result = graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		})
	} else {
		http.Error(w, "No GraphQL query provided", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	r := mux.NewRouter()

	// Handle the root path with a simple message
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the Twitter Clone API!")
	})

	// Handle the /graphql path for GraphQL requests
	r.HandleFunc("/graphql", handleGraphQL).Methods("POST", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(r)

	http.Handle("/", handler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
