package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

type Match struct {
	ID           	int    	`json:"id"`
	Length       	string 	`json:"length"`
	MapName      	string 	`json:"mapName"`
	MatchEnded   	string 	`json:"matchEnded"`
	MatchStarted 	string 	`json:"matchStarted"`
	CT        		int 	`json:"CT"`
	TERRORIST       int 	`json:"TERRORIST"`
}

type Players struct {
	Kills 	int 	`json:"Kills"`
	Died 	int 	`json:"Died"`
	Wins 	int 	`json:"Wins"`
}

type MapInfo struct {
	CT 			int		`json:"ct"`
	TERRORIST 	int		`json:"terrorist"`
	Players    	int    	`json:"Players"`
	KillerNick 	string 	`json:"killerNick"`
}

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

func startMatchInsert(matchID uuid.UUID,mapName string, mapID int) bool {

	db := openDB()
	stmIns, err := db.Prepare("INSERT INTO Matches(id,map_ID,mapName, matchStarted) VALUES (?,?,?, current_timestamp)")
	checkErr(err)

	defer stmIns.Close()

	_ , err = stmIns.Exec(matchID, mapID, mapName)

	if err != nil {
		panic(err.Error())
		return false
	}

	return true
}

func endMatchInsert(ct int,t int, length int, winner string) bool {

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

func getMatchByID(matchID int) ([]Match, error){

	db := openDB()
	var match Match
	var matches []Match

	row, err := db.Prepare("SELECT id,mapName,matchStarted,CT,TERRORIST,matchEnded,matchLength FROM Matches WHERE id= ?")
	checkErr(err)

	rows , err := row.Query(matchID)

	for rows.Next() {
		rows.Scan(&match.ID, &match.MapName, &match.MatchStarted, &match.CT, &match.TERRORIST,&match.MatchEnded, &match.Length)
		matches = append(matches, match)
	}

	db.Close()

	return matches,err
}

func getPlayerInfo(steamID string) ([]Players, error){

	db := openDB()
	var player Players
	var players []Players

	row, err := db.Prepare(`
	SELECT *
	FROM (SELECT COUNT(killerSteamID) as Kills FROM Kills WHERE killerSteamID=?) AS PlayerKills
	JOIN
	(SELECT COUNT(victimSteamID) as Died FROM Kills WHERE victimSteamID=?) AS PlayerDies
	ON 1
	`)
	checkErr(err)

	rows, err := row.Query(steamID, steamID)

	for rows.Next() {
		rows.Scan(&player.Kills, &player.Died)
	}

	row, err = db.Prepare(`
	SELECT DISTINCT COUNT(teamWin) as Wins
	FROM Rounds JOIN Kills
	ON (teamWin=Kills.killerTeam AND Kills.killerSteamID=?) OR (TeamWin=Kills.victimTeam AND Kills.victimSteamID=?) AND Rounds.matchID=Kills.matchID
	GROUP BY(roundID)`)
	checkErr(err)

	rows, err = row.Query(steamID, steamID)

	for rows.Next() {
		rows.Scan(&player.Wins)
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