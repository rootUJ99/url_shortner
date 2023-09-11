package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TinyHandlerBody struct  {
	Url string `json:"url"`
	Expiry int `json:"expiry"`
}


func calculateHash(data string) string {
	sum:=sha256.Sum256([]byte(data))
	cal := fmt.Sprintf("%x", sum)[:9]
	return cal	
}

func sendAsJson(w http.ResponseWriter, requestData any) {
	err:= json.NewEncoder(w).Encode(requestData); if err != nil {
		log.Panic("error while parsing the request data")
	}
}

func parseBody(r *http.Request, body any) {
	err:= json.NewDecoder(r.Body).Decode(body); if err != nil {
		log.Panic("error while parsing the body")
	}
}


func rootRedirector(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.google.com",301)
		// testJson:= test {Name:"holla", Age:27}
		//w.Write(testJson)
		/* json.NewEncoder(w).Encode(testJson) */
		// sendAsJson(w, testJson)

	w.Write([]byte("holla"))
}


func tinyHandler(w http.ResponseWriter, r *http.Request) {
	var body TinyHandlerBody
	//parseBody(r, &body)
	json.NewDecoder(r.Body).Decode(&body)
	hash:=calculateHash(body.Url)
	fmt.Println("hash", hash)
	sendAsJson(w, body)
}

func main(){
	fmt.Println("Hello form go")
	r:= mux.NewRouter()	

	r.HandleFunc("/", rootRedirector).Methods("GET")

	r.HandleFunc("/api/v1/tiny", tinyHandler).Methods("POST")

	http.ListenAndServe(":6969",r)
}
