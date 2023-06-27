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
}