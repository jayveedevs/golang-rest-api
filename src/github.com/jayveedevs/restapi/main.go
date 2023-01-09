package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,POST,DELETE,OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
}

// User Struct (Model)
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
}

var users []User

func getUsers(w http.ResponseWriter, r *http.Request) {

	users = nil

	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sas")

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT UserID, FirstName, LastName, Address, BirthDate, Gender FROM users")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Address, &user.BirthDate, &user.Gender)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, User{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Address: user.Address, BirthDate: user.BirthDate, Gender: user.Gender})

	}

	json.NewEncoder(w).Encode(users)

}

func getUser(w http.ResponseWriter, r *http.Request) {

	users = nil

	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sas")

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var user User

	err = db.QueryRow("SELECT UserID, FirstName, LastName, Address, BirthDate, Gender FROM users where UserID = ?", params["id"]).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Address, &user.BirthDate, &user.Gender)

	// users = append(users, User{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Address: user.Address, BirthDate: user.BirthDate, Gender: user.Gender})

	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sas")

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	insert, err := db.Query("INSERT INTO users (FirstName, LastName, Address, BirthDate, Gender) VALUES ('" + user.FirstName + "','" + user.LastName + "','" + user.Address + "','" + user.BirthDate + "','" + user.Gender + "')")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sas")

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	update, err := db.Query("UPDATE users SET FirstName='" + user.FirstName + "',LastName='" + user.LastName + "',Address='" + user.Address + "',BirthDate='" + user.BirthDate + "',Gender='" + user.Gender + "' WHERE UserID=" + params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer update.Close()

}

func deletUser(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sas")

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	delete, err := db.Query("DELETE FROM users WHERE UserID=" + params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer delete.Close()
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/users/{id}", deletUser).Methods("DELETE", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8082", router))
}
