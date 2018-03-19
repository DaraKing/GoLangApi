package main

import (
	"strings"
	"fmt"
	"regexp"
)

func checkMatchStart(message string) bool {
	 if strings.Index(message, "Match_Start") > 0  {
	 	return true
	 }

	 return false
}

func getMatchIdAndMapname(bodyString string)  {
	info := strings.Split(bodyString, "on")[1]
	id := strings.Split(info, "/")[1]
	mapname := strings.Split(info, "/")[2]
	mapName := strings.TrimSuffix(mapname, "\"")

	fmt.Println(id,mapName)
}

func getKillerAndVictim(bodyString string) (string,string){
	r, _ := regexp.Compile(`"(.*?)"`)
	result := r.FindAllString(bodyString, -1)
	killerPrototype := result[0]
	victimPrototype := result[1]

	return killerPrototype, victimPrototype
}

func getSteamID(message string) string {

	reg, _ := regexp.Compile(`\<(.*?)\>`)
	res := reg.FindAllString(message, -1)

	return res[1]
}

func getTeam(message string) string {
	reg, _ := regexp.Compile(`\<(.*?)\>`)
	res := reg.FindAllString(message, -1)

	return res[2]
}

func getRank(message string) string {
	reg, _ := regexp.Compile(`\<(.*?)\>`)
	res := reg.FindAllString(message, -1)

	return res[0]
}

func getNickName(message string) string {

	return strings.Split(message, getRank(message))[0]

}