package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var initialMap = map[string]string{}

type Error struct {
	ErrorMessage string `json:"error"`
}

func postWebsites(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Link endpoint hit")
	createdMap := map[string][]string{}
	err := json.NewDecoder(r.Body).Decode(&createdMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(createdMap)
	for _, value := range createdMap["websites"] {
		initialMap[value] = "Down"
	}
	fmt.Fprint(w, "Added Succesfully")
	fmt.Println(createdMap)
}

func getWebsites(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get link using Post")
	for key, value := range initialMap {
		fmt.Printf("[%s]=%s\n", key, value)
	}
	jsonStr, err := json.Marshal(initialMap)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(jsonStr))
		fmt.Fprint(w, string(jsonStr))
	}
}
func getWebsitesid(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getwebsite by id endpoint hit")
	param := mux.Vars(r)
	flink := "https://" + param["link"]

	_, present := initialMap[flink]
	if !present {
		errMsg := Error{
			ErrorMessage: "Not Present",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	reponseMp := make(map[string]string)
	reponseMp[flink] = initialMap[flink]
	json.NewEncoder(w).Encode(reponseMp)

}

func getStatus() {
	for {
		for key := range initialMap {
			resp, err := http.Get(key)
			if err != nil {
				fmt.Println("Error occured")
				initialMap[key] = "Down"
				continue
			}
			if resp.StatusCode == 200 {
				fmt.Println("Successful")
				initialMap[key] = "Up"
			}
		}
		time.Sleep(60 * time.Second)
	}
}

func main() {
	go getStatus()
	r := mux.NewRouter()
	r.HandleFunc("/postlink", postWebsites)
	r.HandleFunc("/getlink", getWebsites)
	r.HandleFunc("/getlink/{link}", getWebsitesid)
	fmt.Println("Server on localhost 8081")
	http.ListenAndServe(":8081", r)
}
