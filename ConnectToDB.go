package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func openDB() *sql.DB {
	db, err := sql.Open("mysql", "root:dario123@/CounterStrikeDB")

	if err != nil {
		panic(err.Error())
	}

	return db
}

func checkErr(err error)  {
	if err != nil {
		panic(err.Error())
	}
}

func startMatchInsert(matchID int,mapName string) bool {

	db := openDB()
	stmIns, err := db.Prepare("INSERT INTO Matches(id,mapName) VALUES (?,?)")
	checkErr(err)

	defer stmIns.Close()

	_ , err = stmIns.Exec(matchID,mapName)

	if err != nil {
		panic(err.Error())
		return false
	}

	return true
}