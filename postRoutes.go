package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Create post
func (db DB) CreatePost(user User, w http.ResponseWriter, r *http.Request) {
	post := new(Post)
	// Get body parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}

	// Title and context length must be > 0
	if len(post.Title) == 0 || len(post.Context) == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Please fill all fields")
		return
	}

	// Store post in the database
	query := "INSERT INTO posts (author, title, context) VALUES('" + user.Username + "', '" + post.Title + "', '" + post.Context + "')" 
	rows, err := db.db.Query(query)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	_ = rows

	w.WriteHeader(201)
	json.NewEncoder(w).Encode("Created!")
}

func (db DB) DeletePost(user User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	post := new(Post)

	// Get data from database in order to autherize the request 
	query := "SELECT * FROM posts WHERE id='" + params["id"] + "';"
	rows, err := db.db.Query(query)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	

	// Insert query data to post struct
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.Author, &post.Time_stamp, &post.Title, &post.Context)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}

	// Autherize that the logged in user is the owner of the post
	if post.Author == user.Username {
		// Delete post from database
		query = "DELETE FROM posts WHERE id='" + post.Id + "'"
		rows, err := db.db.Query(query)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode("Internal Server Error")
			return
		}
		_ = rows

		w.WriteHeader(200)
		json.NewEncoder(w).Encode("Post deleted!") 
	} else {
		// Throw error if not autherized
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Unautherized action")
	}
}

// Update
func (db DB) UpdatePost(user User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	post := new(Post)

	// Get data from database in order to autherize the request 
	query := "SELECT * FROM posts WHERE id='" + params["id"] + "';"
	rows, err := db.db.Query(query)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	
	// Insert query data to post struct
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.Author, &post.Time_stamp, &post.Title, &post.Context)
		if err != nil {
			w.WriteHeader(404)
			return
		}
	}

	// Autherize that the logged in user is the owner of the post
	if post.Author == user.Username {
		updatedPost := new(Post)
		// Get the body parameters
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&updatedPost)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode("Internal Server Error")
			return
		}

		// Update changes in database
		if updatedPost.Title != post.Title && len(updatedPost.Title) > 0 {
			query = "UPDATE posts SET title='"+ updatedPost.Title + "' WHERE id='" + post.Id + "'"
			rows, err := db.db.Query(query)
			if err != nil {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode("Internal Server Error")
				return
			}
			_ = rows
		}
		if updatedPost.Context != post.Context && len(updatedPost.Context) > 0  {
			query = "UPDATE posts SET context='"+ updatedPost.Context + "' WHERE id='" + post.Id + "'"
			rows, err := db.db.Query(query)
			if err != nil {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode("Internal Server Error")
				return
			}
			_ = rows
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode("Post updated!") 
	} else {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Unautherized action")
	}
}

// Get all posts paginated
func (db DB) GetPosts(w http.ResponseWriter, r *http.Request) {
	limit := 5
	params := mux.Vars(r)
	// Get page num and parse it as int 
	num, err := strconv.ParseInt(params["page"], 10, 32)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	page := int(num)
	

	// Get data from database with pagination parameters
	query := "SELECT * FROM posts ORDER BY time_stamp DESC LIMIT " + fmt.Sprint(limit) + " OFFSET " + fmt.Sprint((page-1) * limit) 
	rows, err := db.db.Query(query) 
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Map query's data into a list of posts
	var posts []*Post
	for rows.Next() {
		post := new(Post)
		err := rows.Scan(&post.Id, &post.Author, &post.Time_stamp, &post.Title, &post.Context)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		posts = append(posts, post)
	}

	// Get the size of the table
	query = "SELECT COUNT(*) FROM posts"
	rows, err = db.db.Query(query) 
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Insert query's data in count var
	var Count int
	for rows.Next() {
		err := rows.Scan(&Count)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}

	// Get the total pages
	var totalPages int
	if Count % limit > 0 {
		totalPages = Count / limit + 1
	} else {
		totalPages = Count / limit
	}

	// Create response body
	res := map[string]interface{} {
		"posts": posts,
		"currentPage" : page,
		"totalPages": totalPages,
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

// Get post by id
func (db DB) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	post := new(Post)

	// Get post's data from database
	query := "SELECT * FROM posts WHERE id='" +  params["id"] + "'"
	rows, err := db.db.Query(query)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	// Insert query's data in post struct
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.Author, &post.Time_stamp, &post.Title, &post.Context)
		if err != nil {
			w.WriteHeader(404)
			return
		}
	}

	json.NewEncoder(w).Encode(post)
}