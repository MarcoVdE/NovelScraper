package rest_api

import (
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

/**
all of this is based off of:
https://medium.com/@hugo.bjarred/rest-api-with-golang-and-mux-e934f581b8b5

And authentication to be based off of:
https://auth0.com/blog/authentication-in-golang/
and
https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8
*/

//TODO: Test if image works, else base64 encode: encoded := base64.StdEncoding.EncodeToString(content)

var novelList []models.NovelList

func RestAPIEntry(nvlList []models.NovelList) {
	novelList = nvlList
	router := mux.NewRouter()

	//TODO: Include security.

	router.HandleFunc("/novels/", getNovels).Methods("GET")
	router.HandleFunc("/novel-single/", getNovel).
		Queries("title", "{title}").
		Methods("GET")
	router.HandleFunc("/novel-single-chapter/", getNovelChapter).
		Queries("title", "{title}", "chapter", "{chapter}").
		Methods("GET")

	_ = http.ListenAndServe(":8000", router)
}

func getNovels(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(novelList)
}

func getNovel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range novelList {
		if item.Title == params["title"] {
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(&models.NovelList{})
}

func getNovelChapter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range novelList {
		if item.Title == params["title"] {
			chapterRequest, err := strconv.Atoi(params["chapter"])

			if err != nil {
				fmt.Printf("Error: %e on line 117 can't convert given string to int for chapter", err)
				return
			}

			for _, chapterItem := range item.ChapterList {
				//get int version of string request
				if chapterItem.Chapter == chapterRequest {
					_ = json.NewEncoder(w).Encode(chapterItem)
				}
			}
			return
		}
	}
	_ = json.NewEncoder(w).Encode(&models.Chapter{})
}
