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

type TinyGetAllResponse struct {
	Result map[string]string `json:"result"`
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

func redisErr (err error) {
	log.Panic("Error while doing operation with redis", err)
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
		redisErr(err)
		return
	}

	fmt.Println(urlHash)
	http.Redirect(w, r, redisHash,301)
	w.Write([]byte("holla"))
}


func (tCtx TinyCtx) tinyPostHandler(w http.ResponseWriter, r *http.Request) {
	var body TinyHandlerBody
	json.NewDecoder(r.Body).Decode(&body)
	hash:=calculateHash(body.Url)
	tCtx.urls[hash] = body.Url 

	err := tCtx.client.HSet(tCtx.ctx, "urlHash", body.Url, hash).Err()
	if err != nil {
		redisErr(err)
	}

	fmt.Println(tCtx.urls)
	sendAsJson(w, body)
}

func (tCtx TinyCtx) tinyGetHandler(w http.ResponseWriter, r *http.Request) {
	url:=r.URL.Query().Get("url")
	resultHash, err := tCtx.client.HGet(tCtx.ctx, "urlHash", url).Result()
	if err != nil {
		redisErr(err)
		return
	}

	fmt.Println(url, resultHash, "here we go is there any space")
	
	type res struct {
		Result map[string]string `json:"result"`
	}
	response:=res {
		map[string]string{url:resultHash },
	}
	sendAsJson(w, response)
}

func (tCtx TinyCtx) tinyGetAllHandler(w http.ResponseWriter, r *http.Request) {
	resultList, err := tCtx.client.HGetAll(tCtx.ctx, "urlHash").Result()
	if err != nil {
		redisErr(err)
		return
	}

	fmt.Println(resultList)
	
	response:=TinyGetAllResponse {
		resultList,
	}
	sendAsJson(w, response)
}

func (tCtx TinyCtx) tinyDelHandler(w http.ResponseWriter, r *http.Request) {
	url:=r.URL.Query().Get("url")
	err :=tCtx.client.HDel(tCtx.ctx, "urlHash", url).Err()
	if err != nil {
		redisErr(err)
		return
	}
	fmt.Println(url)
	type res struct {
		Message string `json:"message"`
	}
	response:=res {
		Message: fmt.Sprintf("%v has been deleted from the record", url),
	}
	sendAsJson(w, response)
}


func main(){
	fmt.Println("Hello form golang")
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

	r.HandleFunc("/api/v1/tiny", tCtx.tinyPostHandler).Methods("POST", "PUT")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyGetHandler).Methods("GET")

	r.HandleFunc("/api/v1/tiny/all", tCtx.tinyGetAllHandler).Methods("GET")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyDelHandler).Methods("DELETE")

	http.ListenAndServe(":6969",r)
}
