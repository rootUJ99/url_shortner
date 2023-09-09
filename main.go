package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type test struct  {
	Name string `json:"name"`
	Age int `json:"age"`
}

func sendAsJson(w http.ResponseWriter, requestData any) {
	err:= json.NewEncoder(w).Encode(requestData); if err != nil {
		log.Panic("error while parsing the request data")
	}
}

func main(){
	fmt.Println("Hello form go")
	r:= mux.NewRouter()	
	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		// testJson:= test {Name:"holla", Age:27}
		//w.Write(testJson)
		/* json.NewEncoder(w).Encode(testJson) */
		// sendAsJson(w, testJson)
		http.Redirect(w, r, "https://www.google.com",301)
	}).Methods("GET")
	http.ListenAndServe(":6969",r)
}
