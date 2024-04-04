package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
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

type stuDetails struct {
	cid     int
	stname  string
	stemail string
	stphno  string
}

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
		fmt.Println("Preparation Error")
		log.Fatal(err)
	}
	defer insStmt.Close()

	_, err = insStmt.Exec(cid, stname, stemail, stphno)
	if err != nil {
		fmt.Println("Execution Error")
		http.Error(w, "Cid exists", http.StatusConflict)
		return
	}
    
    pwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Path Error")
        return
    }

    pwd = path.Join(pwd, "/index.html")
    fmt.Println(pwd)

	db.Close()
	http.Redirect(w, r, pwd, http.StatusSeeOther)
}

func tableViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // or a specific domain
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, hx-request")

	db, err := sql.Open("sqlite3", "trial.db")
	if err != nil {
		fmt.Fprintln(w, "Sqlite3 error")
		return
	}
	defer db.Close()

	err = r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, "Form Error")
		return
	}

	cid, err := strconv.Atoi(r.FormValue("cid"))
	if err != nil {
		fmt.Println(err, "Conversion Error")
		return
	}

	stmt, err := db.Prepare("SELECT * FROM studata WHERE cid=?")
	if err != nil {
		fmt.Fprintln(w, "Prepare error")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(cid)
	if err != nil {
		fmt.Println("Query error")
		log.Fatalln(err)
		return
	}
	defer rows.Close()

	tableSlice := []string{}
	for rows.Next() {
		var newdata stuDetails
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
	db.Close()
}

func main() {
	fmt.Println("Server Started . . .")

	http.HandleFunc("/submit", dataHandler)

	http.HandleFunc("/getres", tableViewer)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
