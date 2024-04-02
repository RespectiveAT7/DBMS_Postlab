package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	file        string = "trial.db"
	createTable string = `
CREATE TABLE IF NOT EXISTS studata (
	cid INT PRIMARY KEY,
	stname TEXT,
	stemail TEXT,
	stphno TEXT
)
`
	tableName string = "studata"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("Only POST allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Parsing Error")
		return
	}

	cid, err := strconv.Atoi(r.FormValue("cid"))
	if err != nil {
		fmt.Println("Cannot insert into cid")
		log.Fatal(err)
		return
	}
	stname := r.FormValue("stname")
	stemail := r.FormValue("stemail")
	stphno := r.FormValue("stphno")

	
	db, err := sql.Open("sqlite3", "trial.db")
	if err != nil {
		fmt.Println("Sqlite Error: cannot open")
	}
	defer db.Close()
	
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Create table error")
		return
	}
	
	insStmt, err := db.Prepare("INSERT INTO studata(cid, stname, stemail, stphno) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insStmt.Close()
	
	_, err = insStmt.Exec(cid, stname, stemail, stphno)
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func tableViewer(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "trial.db")
	if err != nil {
		fmt.Println("Sqlite3 error")
		return
	}

	
}

func main() {
	fmt.Println("Server Started . . .")

	http.HandleFunc("/submit", dataHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
