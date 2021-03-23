package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
)

// Sign user up
func (db DB) UserSignUp(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if len(user.Username) < 6 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Username must be at least 8 characters long")
		return
	}

	if len(user.Password) < 6 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Password must be at least 6 characters long")
		return
	}

	// Hash password before storing into database
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}

	qeury := "INSERT INTO users (username, password) VALUES ('" + user.Username + "', '" + string(hashedPassword) + "')"

	// Storing user in database
	rows, err := db.db.Query(qeury)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Username is already taken")
		return
	}
	_ = rows

	w.WriteHeader(201)
	json.NewEncoder(w).Encode("Signed up!")
}


// Log user in
func (db DB) UserLogin(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	// Store inserted password as bytes to compare with hash
	inputPassword := []byte(user.Password)

	// Get user's data
	query := "SELECT * FROM users WHERE username='" + user.Username + "';"
	rows, err := db.db.Query(query)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	// Insert queries data into user's struct
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode("Internal Server Error")
			return
		}
	}

	// Convert password to bytes
	storedPassword := []byte(user.Password)
	// Validate password
	err = bcrypt.CompareHashAndPassword(storedPassword, inputPassword)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	// Initialiaze jwt claims map
	jwtClaims := jwt.MapClaims{}
	jwtClaims["id"] = user.Id
	// Hash token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString([]byte(jwtSecret))
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}

	// Create cookie
	cookie := http.Cookie{
		Name: "auth_token",
		Value: token,	
		MaxAge: 900000,
		// HttpOnly: true,
		Path: "/",
	}
	// Set cookie 
	http.SetCookie(w, &cookie)
	w.WriteHeader(200)
}


// Log user out
func UserLogout(w http.ResponseWriter, r *http.Request) {
	// Create expired cookie
	cookie := http.Cookie{
		Name:   "auth_token",
		MaxAge: -1,
		HttpOnly: true,
		Path: "/",
	}

	// Set auth cookie to the expired in order to be removed by the brower
	http.SetCookie(w, &cookie)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("Logged out!")
}

// Define type of authenticated user's function
type userFunc = func(user User, w http.ResponseWriter, r *http.Request)

// Authenticate and autherize user middleware
func (db DB) AuthenticateUser(next userFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get auth cookie
		auth_cookie, err := r.Cookie("auth_token")
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode("Please")
			return
		}
		jwtClaims := jwt.MapClaims{}
		// Validate auth cookie value
		token, err := jwt.ParseWithClaims(auth_cookie.Value, jwtClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode("Unautherized action")
			return
		}
		_ = token

		// Get id of user
		var id string
		for key, value := range jwtClaims {
			if key == "id" {
				id = value.(string)
			}
		}

		// Query user
		query := "SELECT * FROM users WHERE id='" + id + "'"
		rows, err := db.db.Query(query)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode("Internal Server Error")
			return
		}

		// Insert query's data to user's struct
		user := new(User)
		for rows.Next() {
			err := rows.Scan(&user.Id, &user.Username, &user.Password)
			if err != nil {
				w.WriteHeader(404)
				json.NewEncoder(w).Encode("User not found")
				return
			}
		}

		// Pass parameters to authenticated user's function
		next(*user, w, r)
	})
}

// Get user's profile 
func (db DB) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get url passed parameters
	params := mux.Vars(r)
	user := new(User)

	limit := 5
	// Get page num and parse it as int 
	num, err := strconv.ParseInt(params["page"], 10, 32)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	page := int(num)

	// Query user's data
	query := "SELECT * FROM users WHERE username='" + params["username"] + "'"
	rows, err := db.db.Query(query)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("User not found")
		return
	}

	// Insert query's data to user's struct
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode("Internal Server Error")
			return
		}
	}

	// Query user's post data
	query = "SELECT * FROM posts WHERE author='" + user.Username + "' ORDER BY time_stamp DESC LIMIT " + fmt.Sprint(limit) + " OFFSET " + fmt.Sprint((page-1) * limit) 
	rows, err = db.db.Query(query)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}

	// Insert query's data to posts's structs
	var posts []*Post
	for rows.Next() {
		post := new(Post)
		err := rows.Scan(&post.Id, &post.Author, &post.Time_stamp, &post.Title, &post.Context)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		posts = append(posts, post)
	}

	// Get the size of the table
	query = "SELECT COUNT(*) FROM posts WHERE author='" + user.Username + "'"
	rows, err = db.db.Query(query) 
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}

	// Insert query's data in count var
	var Count int
	for rows.Next() {
		err := rows.Scan(&Count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode("Internal Server Error4")
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
	res := map[string]interface{}{
		"username": user.Username,
		"posts":    posts,
		"totalPages": totalPages,
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

// Get logged in user's data
func LoggedUserProfile(user User, w http.ResponseWriter, r *http.Request) { 
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(user)
}