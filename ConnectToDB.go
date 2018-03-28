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

func startMatchInsert(mapName string) bool {

	db := openDB()
	stmIns, err := db.Prepare("INSERT INTO Matches(id,mapName, matchStarted) VALUES (?,?, current_timestamp)")
	checkErr(err)

	defer stmIns.Close()

	_ , err = stmIns.Exec(matchID,mapName)

	if err != nil {
		panic(err.Error())
		return false
	}

	return true
}

func endMatchInsert(ct int,t int, length string, winner string) bool {

	db := openDB()
	stmt, err := db.Prepare("UPDATE Matches SET CT=?, TERRORIST=?, matchEnded=current_timestamp, matchLength=?, winner=? WHERE id=?")
	checkErr(err)

	result, err := stmt.Exec(ct,t,length,winner,matchID)

	defer stmt.Close()

	var numResults, _ = result.RowsAffected()

	if err !=nil || numResults != 1 {
		return false
	}

	return true
}

func insertInKillsTable(killerNick string, killerSteamID string, killerTeam string, headshot bool, victimNick string, victimSteamID string, victimTeam string, weapon string) bool {

	db := openDB()
	stmIns, err := db.Prepare("INSERT INTO Kills(matchID, killerNick, killerSteamID, killerTeam, headshot, victimNick, victimSteamID, victimTeam, weapon) VALUES (?,?,?,?,?,?,?,?,?)")
	checkErr(err)

	defer stmIns.Close()

	_, err = stmIns.Exec(matchID,killerNick,killerSteamID,killerTeam,headshot,victimNick,victimSteamID,victimTeam,weapon)

	if err != nil {
		panic(err.Error())
		return false
	}

	return true
}

func insertRound(teamWin string, currentScore string) bool{

	db := openDB()
	stmIns, err := db.Prepare("INSERT INTO Rounds (matchID, teamWin, currentScore) VALUES (?,?,?)")
	checkErr(err)

	defer stmIns.Close()

	_, err = stmIns.Exec(matchID, teamWin, currentScore)

	if err !=  nil {
		panic (err.Error())
		return false
	}

	return true
}

type Match struct {
	ID           int    `json:"id"`
	Length       string `json:"length"`
	MapName      string `json:"mapName"`
	MatchEnded   string `json:"matchEnded"`
	MatchStarted string `json:"matchStarted"`
	Score        string `json:"score"`
}

type Players struct {
	Kills int `json:"Kills"`
	Died int `json:"Died"`
}

type MapInfo struct {
	CT 			int		`json:"ct"`
	TERRORIST 	int		`json:"terrorist"`
	Players    	int    	`json:"Players"`
	KillerNick 	string 	`json:"killerNick"`
}

func getMatchByID(matchID int) ([]Match, error){

	db := openDB()
	var match Match
	var matches []Match

	row, err := db.Prepare("SELECT id,mapName,matchStarted,score,matchEnded,length FROM Matches WHERE id= ?")
	checkErr(err)

	rows , err := row.Query(matchID)

	for rows.Next() {
		rows.Scan(&match.ID, &match.MapName, &match.MatchStarted, &match.Score ,&match.MatchEnded, &match.Length)
		matches = append(matches, match)
	}

	db.Close()

	return matches,err
}

func getPlayerInfo(steamID string) ([]Players, error){

	db := openDB()
	var player Players
	var players []Players

	row, err := db.Prepare("SELECT COUNT(killerSteamID) as Kills FROM Kills WHERE killerSteamID=?")
	checkErr(err)

	rows, err := row.Query(steamID)

	for rows.Next() {
		rows.Scan(&player.Kills)
	}

	row, err = db.Prepare("SELECT COUNT(victimSteamID) as Died FROM Kills WHERE victimSteamID=?")
	checkErr(err)

	rows, err = row.Query(steamID)

	for rows.Next() {
		rows.Scan(&player.Died)
		players = append(players, player)
	}

	db.Close()

	return players, err
}

func getMapInfo(name string) ([]MapInfo, error){

	db := openDB()
	var mapInfo MapInfo
	var info []MapInfo

	row, err := db.Prepare("SELECT CT,TERRORIST FROM Matches WHERE mapName=?")
	checkErr(err)
	rows, err := row.Query(name)

	for rows.Next() {
		rows.Scan(&mapInfo.CT, &mapInfo.TERRORIST)
	}

	row, err = db.Prepare("SELECT Kills.killerNick, COUNT(*) as Players FROM Kills,Matches WHERE Matches.mapName=? AND Matches.id=Kills.matchID GROUP BY Kills.killerNick")
	checkErr(err)
	rows, err = row.Query(name)

	for rows.Next() {
		rows.Scan(&mapInfo.KillerNick, &mapInfo.Players)
		info = append(info, mapInfo)
	}

	db.Close()

	return info,err
}