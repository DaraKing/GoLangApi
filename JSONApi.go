package main

import (
	"net/http"
	"io/ioutil"
)

func optionsRequest(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func getMessage(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)

	bodyBytes, _ := ioutil.ReadAll(r.Body)

	bodyString := string(bodyBytes)

	if checkMatchStart(bodyString) {
		getMatchIdAndMapname(bodyString)
	}

	if checkIsKill(bodyString) {
		getInfoAboutKill(bodyString)
	}

	if checkRoundOver(bodyString) {
		getRoundInfo(bodyString)
	}

	if isGameOver(bodyString) {
		getGameStats(bodyString)
	}
}