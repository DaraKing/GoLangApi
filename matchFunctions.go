package main

import (
	"strings"
	"regexp"
	"strconv"
	"github.com/satori/go.uuid"
)

var matchID uuid.UUID

func getRegExParams(regEx, inputString string) map[string]string {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(inputString)

	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func checkMatchStart(message string) bool {
	return strings.Index(message, `"Match_Start"`) > 0
}

func checkIsKill(message string) bool {
	return strings.Index(message, "killed") > 0
}

func checkRoundOver(message string) bool {
	return strings.Index(message, "SFUI_Notice") > 0
}

func checkIsHeadshot(message string) bool {
	return message == "(headshot)"
}

func isGameOver(message string) bool {
	r, _ := regexp.Compile(`Game Over:`)
	res := r.FindAllString(message, -1)

	if res != nil {
		return true
	}

	return false
	//return strings.Index(message, `Game Over: `) > 0
}

func getRoundInfo(bodyString string) bool {

	CSRoundRegEx := `^Team "(?P<teamWinner>.+)" triggered ".+" \(CT "(?P<ctScore>.+)"\) \(T "(?P<tScore>.+)"\)$`
	CSVars := getRegExParams(CSRoundRegEx, bodyString)

	score := CSVars["ctScore"]+":"+CSVars["tScore"]

	return insertRound(CSVars["teamWinner"], score)
}

func getInfoAboutKill(bodyString string) bool {
	CSKillRegEx := `^"(?P<userOneName>.+)<\d+><(?P<userOneSteamId>.+)><(?P<userOneTeam>.+)>" \[.+\] killed "(?P<userTwoName>.+)<\d+><(?P<userTwoSteamId>.+)><(?P<userTwoTeam>.+)>" \[.+\] with "(?P<weaponName>.+)" ?(?P<isHeadshot>\(.*\))?$`
	CSVars := getRegExParams(CSKillRegEx, bodyString)
	return insertInKillsTable(CSVars["userOneName"],CSVars["userOneSteamId"],CSVars["userOneTeam"],checkIsHeadshot(CSVars["isHeadshot"]),
		CSVars["userTwoName"], CSVars["userTwoSteamId"], CSVars["userTwoTeam"], CSVars["weaponName"])

}

func checkWinner(ct int, t int) string {
	if ct > t {
		return "CT"
	}

	return "TERRORIST"
}

func getGameStats(bodyString string) bool {

	CSMatchEndRegEx := `^Game Over: casual \d+ workshop/(?P<mapID>\d+)/(?P<mapName>[a-zA-Z]+_[a-zA-Z0-9]+) score (?P<ct>\d+):(?P<t>\d+) after (?P<minutes>\d+) min$`
	CSVars := getRegExParams(CSMatchEndRegEx, bodyString)

	ct, _ := strconv.Atoi(CSVars["ct"])
	t, _ := strconv.Atoi(CSVars["t"])
	length, _ := strconv.Atoi(CSVars["minutes"])

	return endMatchInsert(ct, t, length, checkWinner(ct,t))
}

func getMatchIdAndMapname(bodyString string)  bool{
	CSMatchRegEx := `World triggered "Match_Start" on "workshop/(?P<mapId>[0-9]+)/(?P<mapName>[a-zA-Z]+_[a-zA-Z0-9]+)"`
	CSVars := getRegExParams(CSMatchRegEx, bodyString)
	mapID, _ := strconv.Atoi(CSVars["mapId"])
	matchID = uuid.Must(uuid.NewV4())
	return startMatchInsert(matchID, CSVars["mapName"], mapID)

}

//func getKillerAndVictimAndWeapon(bodyString string) (string,string,string){
//	r, _ := regexp.Compile(`"(.*?)"`)
//	result := r.FindAllString(bodyString, -1)
//
//	return result[0], result[1], removeQuotes(result[2])
//}
//
//func removeChar(char string) string{
//	reg, _ := regexp.Compile(`\<|\>`)
//	res := reg.ReplaceAllString(char,``)
//
//	return res
//}
//
//func removeQuotes(char string) string {
//	reg, _ := regexp.Compile(`\"`)
//	res := reg.ReplaceAllString(char ,``)
//
//	return res
//}
//
//func removeSlash(char string) string {
//	reg, _ := regexp.Compile(`\/`)
//	res := reg.ReplaceAllString(char ,``)
//
//	return res
//}
//
//func split(message string) []string {
//	reg, _ := regexp.Compile(`\<(.*?)\>`)
//	res := reg.FindAllString(message, -1)
//
//	return res
//}
//
//func getSteamID(message string) string {
//	res := split(message)
//
//	return removeChar(res[1])
//}
//
//func getTeam(message string) string {
//	res := split(message)
//
//	return removeChar(res[2])
//}
//
//func getNickName(message string) string {
//
//	reg, _ := regexp.Compile(`\<(.*?)\>|\"`)
//	res := reg.ReplaceAllString(message,``)
//
//	return res
//
//}