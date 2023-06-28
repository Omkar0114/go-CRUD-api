package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

type User struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"Email"`
}

func main(){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT);")
	if err != nil {
		log.Fatal(err)
	}

// create router
router := mux.NewRouter()
router.HandleFunc("/users", getUsers(db)).Methods("GET")
router.HandleFunc("/users/{id}", getUserbyId(db)).Methods("GET")
router.HandleFunc("/users", createUser(db)).Methods("POST")
router.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
router.HandleFunc("/users/{id}", deleteUser(db)).Methods("DELETE")

// start server
log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			user := User{}
			err := rows.Scan(&user.Id, &user.Name, &user.Email)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
		json.NewEncoder(w).Encode(users)
	}
} 

// get user by id
func getUserbyId(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id := vars["id"]
		row := db.QueryRow("SELECT * FROM users WHERE id=$1;", id)
		user := User{}
		err := row.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(user)
	}
}

// create user
func createUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2);", user.Name, user.Email)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(user)
	}
}

// update user
func updateUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id := vars["id"]
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		_, err := db.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3;", user.Name, user.Email, id)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(user)
	}
}

// delete user
func deleteUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		id := vars["id"]
		_, err := db.Exec("DELETE FROM users WHERE id=$1;", id)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(id)
		fmt.Fprintf(w, "User deleted")
	}
}