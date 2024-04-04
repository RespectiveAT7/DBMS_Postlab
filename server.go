package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

var (
	db      *sql.DB
	insStmt *sql.Stmt
	selStmt *sql.Stmt
	allStmt *sql.Stmt
)

type stuDetails struct {
	cid     int
	stname  string
	stemail string
	stphno  string
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "trial.db")
	if err != nil {
		fmt.Println("Sqlite Error: cannot open")
		return
	}

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Create table error")
		return
	}

	insStmt, err = db.Prepare("INSERT INTO studata(cid, stname, stemail, stphno) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println("ins Preparation Error")
		return
	}

	selStmt, err = db.Prepare("SELECT * FROM studata WHERE cid=?")
	if err != nil {
		fmt.Println("sel Prepare error")
		return
	}

	allStmt, err = db.Prepare("SELECT * FROM studata")
	if err != nil {
		fmt.Println("all Prepare error")
		return
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("Only POST Method allowed")
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

	_, err = insStmt.Exec(cid, stname, stemail, stphno)
	if err != nil {
		fmt.Println("Execution Error")
		http.Error(w, "Cid exists", http.StatusConflict)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func tableViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // or a specific domain
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, hx-request")

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, "Form Error")
		return
	}

	cid, err := strconv.Atoi(r.FormValue("cid"))
	if err != nil {
		fmt.Println(err, "Conversion Error")
		return
	}

	rows, err := selStmt.Query(cid)
	if err != nil {
		fmt.Println("Query error")
		log.Fatalln(err)
		return
	}
	defer rows.Close()

	tableSlice := []string{}
	var newdata stuDetails
	for rows.Next() {
		err = rows.Scan(&newdata.cid, &newdata.stname, &newdata.stemail, &newdata.stphno)
		if err != nil {
			fmt.Println("row scan error")
			log.Fatalln(err)
			return
		}
		newRow := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", newdata.cid, newdata.stname, newdata.stemail, newdata.stphno)
		tableSlice = append(tableSlice, newRow)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Rows iteration error", err)
		return
	}

	joinedString := strings.Join(tableSlice, " ")
	fmt.Println(joinedString)
	fmt.Fprint(w, joinedString)
}

func getAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // or a specific domain
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, hx-request")

	rows, err := allStmt.Query()
	if err != nil {
		fmt.Println("Query error")
		log.Fatalln(err)
		return
	}
	defer rows.Close()

	tableSlice := make([]string, 0)
	var newdata stuDetails
	for rows.Next() {
		err = rows.Scan(&newdata.cid, &newdata.stname, &newdata.stemail, &newdata.stphno)
		if err != nil {
			fmt.Println("row scan error")
			log.Fatalln(err)
			return
		}
		newRow := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", newdata.cid, newdata.stname, newdata.stemail, newdata.stphno)
		tableSlice = append(tableSlice, newRow)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Rows iteration error", err)
		return
	}

	joinedString := strings.Join(tableSlice, " ")
	fmt.Fprint(w, joinedString)
}

func main() {
	defer insStmt.Close()
	defer selStmt.Close()
	defer allStmt.Close()
	defer db.Close()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Println("Server Started . . .")
    fmt.Println("Hosting at: http://localhost:8081")

	http.HandleFunc("/submit", dataHandler)

	http.HandleFunc("/getres", tableViewer)

	http.HandleFunc("/getall", getAll)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
