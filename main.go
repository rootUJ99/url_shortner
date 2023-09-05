package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main(){
	fmt.Println("Hello form go")
	r:= mux.NewRouter()	
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is gorilla\n"))
	})
	http.ListenAndServe(":6969",r)
}
