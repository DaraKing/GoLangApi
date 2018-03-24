package main

import (
	"strings"
	"fmt"
	"regexp"
	"strconv"
)

var matchID int

func checkMatchStart(message string) bool {
	 if strings.Index(message, "Match_Start") > 0  {
	 	return true
	 }

	 return false
}

func checkIsKill(message string) bool {
	if strings.Index(message, "killed") > 0  {
		return true
	}

	return false
}

func checkRoundOver(message string) bool {
	if strings.Index(message, "SFUI_Notice") > 0 {
		return true
	}

	return false
}

func checkIsHeadshot(message string) bool {
	if strings.Index(message, "headshot") > 0 {
		return true
	}

	return false
}

func isGameOver(message string) bool {
	r, _ := regexp.Compile(`Game Over:`)
	res := r.FindAllString(message, -1)

	if res != nil {
		return true
	}

	return false
}

func getRoundInfo(bodyString string) {
	r, _ := regexp.Compile(`"(.*?)"`)
	teamWin := r.FindAllString(bodyString, -1)

	regexNum, _ := regexp.Compile(`"([0-9])"`)
	RoundResult := regexNum.FindAllString(bodyString, -1)

	fmt.Println("CT: " +removeQuotes(RoundResult[0]))
	fmt.Println("T: " +removeQuotes(RoundResult[1]))

	fmt.Println(removeQuotes(teamWin[0]))
}

func getInfoAboutKill(bodyString string) {
	killer, victim, weapon := getKillerAndVictimAndWeapon(bodyString)

	fmt.Println("Headshot: " ,checkIsHeadshot(bodyString))

	killerNick := getNickName(killer)
	killerSteamID := getSteamID(killer)
	killerTeam := getTeam(killer)

	victimNick := getNickName(victim)
	victimSteamID := getSteamID(victim)
	victimTeam := getTeam(victim)

	fmt.Println(matchID)
	fmt.Println(killerNick + " - " + killerSteamID + " - " + killerTeam + " - " +weapon)
	fmt.Println(victimNick + " - " + victimSteamID + " - " + victimTeam)
}

func getGameStats(bodyString string) {

	getMatchIdAndMapname(bodyString)

	//TIME
	r, _ := regexp.Compile(`[0-9]+ min`)
	time := r.FindAllString(bodyString, -1)
	fmt.Println("Time is: " +time[0])

	//SCORE
	r, _ = regexp.Compile(`[0-9]+:[0-9]+`)
	score := r.FindAllString(bodyString, -1)
	fmt.Println("Score is: " +score[0])

}

func getMatchIdAndMapname(bodyString string)  string{
	r, _ := regexp.Compile(`/[0-9]*/`)
	id := r.FindAllString(bodyString, -1)
	fmt.Println("Id: " +removeSlash(id[0]))

	matchID, _  = strconv.Atoi(removeSlash(id[0]))

	//MAP
	r, _ = regexp.Compile(`/[a-zA-Z]+_[a-zA-Z|0-9]+`)
	mapName := r.FindAllString(bodyString, -1)
	return removeSlash(mapName[0])

}

func getKillerAndVictimAndWeapon(bodyString string) (string,string,string){
	r, _ := regexp.Compile(`"(.*?)"`)
	result := r.FindAllString(bodyString, -1)

	return result[0], result[1], removeQuotes(result[2])
}

func removeChar(char string) string{
	reg, _ := regexp.Compile(`\<|\>`)
	res := reg.ReplaceAllString(char,``)

	return res
}

func removeQuotes(char string) string {
	reg, _ := regexp.Compile(`\"`)
	res := reg.ReplaceAllString(char ,``)

	return res
}

func removeSlash(char string) string {
	reg, _ := regexp.Compile(`\/`)
	res := reg.ReplaceAllString(char ,``)

	return res
}

func getSteamID(message string) string {

	reg, _ := regexp.Compile(`\<(.*?)\>`)
	res := reg.FindAllString(message, -1)

	return removeChar(res[1])
}

func getTeam(message string) string {
	reg, _ := regexp.Compile(`\<(.*?)\>`)
	res := reg.FindAllString(message, -1)

	return removeChar(res[2])
}

func getNickName(message string) string {

	reg, _ := regexp.Compile(`\<(.*?)\>|\"`)
	res := reg.ReplaceAllString(message,``)

	return res

}