package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"log"
	"encoding/json"
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
		mapName := getMatchIdAndMapname(bodyString)
		status := startMatchInsert(mapName)

		if status {
			w.WriteHeader(http.StatusOK)
		} else {
			fmt.Println(http.StatusInternalServerError)
		}
	}

	if checkIsKill(bodyString) {
		status := getInfoAboutKill(bodyString)

		if status {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if checkRoundOver(bodyString) {
		status := getRoundInfo(bodyString)

		if status {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if isGameOver(bodyString) {
		status := getGameStats(bodyString)

		if status {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getMatch(w http.ResponseWriter, r *http.Request)  {
	setupResponse(&w, r)

	vars := mux.Vars(r)
	key := vars["id"]
	id,err := strconv.Atoi(key)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	matches, error := getMatchByID(id)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error while querying")
	}

	json.NewEncoder(w).Encode(matches)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	vars := mux.Vars(r)
	key := vars["steamID"]

	player, err := getPlayerInfo(key)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(player)
}

func mapInfo(w http.ResponseWriter, r *http.Request)  {
	setupResponse(&w, r)

	vars := mux.Vars(r)
	key := vars["name"]

	mapInfo, err := getMapInfo(key)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(mapInfo)
}