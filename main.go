package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
)

const (
	dbURI     = ""
	jwtSecret = ""
)

func main() {
	// DATABASE
	var err error

	db := new(DB)
	// Connect to database
	db.db, err = sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}
	defer db.db.Close()

	err = db.db.Ping()
	if err != nil {
		panic(err)
	}
	println("Database connected!")

	// ROUTES
	router := mux.NewRouter()

	router.HandleFunc("/user/signup", db.UserSignUp).Methods("POST")
	router.HandleFunc("/user/login", db.UserLogin).Methods("POST")
	router.HandleFunc("/user/logout", UserLogout).Methods("DELETE")
	router.Handle("/user/me", db.AuthenticateUser(LoggedUserProfile)).Methods("GET")
	router.Handle("/create/post", db.AuthenticateUser(db.CreatePost)).Methods("POST")
	router.Handle("/delete/post/{id}", db.AuthenticateUser(db.DeletePost)).Methods("POST")
	router.Handle("/update/post/{id}", db.AuthenticateUser(db.UpdatePost)).Methods("POST")
	router.HandleFunc("/users/{username}/{page}", db.GetUserProfile).Methods("GET")
	router.HandleFunc("/posts/pages/{page}", db.GetPosts).Methods("GET")
	router.HandleFunc("/post/{id}", db.GetPost).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
		// Debug: true,
	})

	handler := c.Handler(router)

	// SERVER
	println("Listening on port :8080...")
	http.ListenAndServe(":8080", handler)
}