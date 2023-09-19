package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"context"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type TinyHandlerBody struct  {
	Url string `json:"url"`
	Expiry int `json:"expiry"`
}


type TinyCtx struct {
	urls map[string]string
	client *redis.Client
	ctx context.Context 
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


func (tCtx TinyCtx) rootRedirector(w http.ResponseWriter, r *http.Request) {
	hash:=mux.Vars(r)["hash"]
	// if hash == nil {
	// 	log.Panic("pass the hash")	
	// 	return
	// }
	urlHash, ok := tCtx.urls[hash]
	if !ok {
		log.Panic("url does not exist")
		return
	}
	redisHash ,err := tCtx.client.Get(tCtx.ctx, hash).Result()
	if err != nil {
	    panic(err)
	}

	fmt.Println(urlHash)
	http.Redirect(w, r, redisHash,301)
		// testJson:= test {Name:"holla", Age:27}
		//w.Write(testJson)
		/* json.NewEncoder(w).Encode(testJson) */
		// sendAsJson(w, testJson)

	w.Write([]byte("holla"))
}


func (tCtx TinyCtx) tinyHandler(w http.ResponseWriter, r *http.Request) {
	var body TinyHandlerBody
	//parseBody(r, &body)
	json.NewDecoder(r.Body).Decode(&body)
	hash:=calculateHash(body.Url)
	fmt.Println("hash", hash)
	// Append(tinyList, hash) 
	tCtx.urls[hash] = body.Url 
	err := tCtx.client.Set(tCtx.ctx, hash, body.Url, 0).Err()
	if err != nil {
	    panic(err)
	}
	fmt.Println(tCtx.urls)
	sendAsJson(w, body)
}

func main(){
	fmt.Println("Hello form go")
	r:= mux.NewRouter()	
	client := redis.NewClient(&redis.Options{

        Addr:	  "redis:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })
	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
	    panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
	    panic(err)
	}
	fmt.Println("foo", val)
	tCtx :=TinyCtx{
		urls: make(map[string]string),
		client: client,
		ctx: ctx,
	}
	r.HandleFunc("/{hash}",tCtx.rootRedirector).Methods("GET")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyHandler).Methods("POST")

	http.ListenAndServe(":6969",r)
}
