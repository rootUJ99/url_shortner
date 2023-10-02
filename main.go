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
	"github.com/rootuj99/url_shortner/cmd/cli"
)

type TinyHandlerPostBody struct  {
	Url string `json:"url"`
	Expiry int `json:"expiry"`
}

type TinyHandlerUpdateBody struct  {
	OldUrl string `json:"oldurl"`
	Url string `json:"url"`
	Expiry int `json:"expiry"`
}

type TinyGetAllResponse struct {
	Result map[string]string `json:"result"`
} 


type TinyCtx struct {
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

func sendErrJson(w http.ResponseWriter, message string) {
	type errMessage struct {
		Message string `json:"message"`
	} 
	errorRes:= errMessage{
		Message: message,
	}
	w.WriteHeader(http.StatusBadRequest)
	err:= json.NewEncoder(w).Encode(errorRes); if err != nil {
		log.Panic("error while parsing the request data")
	}
}


func redisErr (err error) {
	log.Panic("Error while doing operation with redis", err)
}

func (tCtx TinyCtx) rootRedirector(w http.ResponseWriter, r *http.Request) {
	hash:=mux.Vars(r)["hash"]

	redisHash ,err := tCtx.client.HGet(tCtx.ctx, "urlHash", hash).Result()

	fmt.Println("this should be the hash", redisHash)
	if err != nil {
		sendErrJson(w, fmt.Sprintf("%v%v url does not exist!", r.Host, r.URL))
		return
	}

	http.Redirect(w, r, redisHash,301)
	w.Write([]byte(fmt.Sprintf("redirecting to %v", redisHash)))
}


func (tCtx TinyCtx) tinyPostHandler(w http.ResponseWriter, r *http.Request) {
	var body TinyHandlerPostBody 
	json.NewDecoder(r.Body).Decode(&body)
	hash:=calculateHash(body.Url)
	
	err := tCtx.client.HSet(tCtx.ctx, "urlHash", hash, body.Url).Err()
	if err != nil {
		redisErr(err)
		return
	}

	type res struct {
		Message string `json:"message"`
	}
	response:=res {
		Message: fmt.Sprintf("%v has been added in the record", body.Url),
	}
	sendAsJson(w, response)
}

func (tCtx TinyCtx) tinyPutHandler(w http.ResponseWriter, r *http.Request) {
	var body TinyHandlerUpdateBody 
	json.NewDecoder(r.Body).Decode(&body)
	hash:=calculateHash(body.Url)
	
	err := tCtx.client.HSet(tCtx.ctx, "urlHash", hash, body.Url).Err()
	if err != nil {
		redisErr(err)
		return
	}

	err = tCtx.client.HDel(tCtx.ctx, "urlHash", body.OldUrl).Err()
	if err != nil {
		redisErr(err)
		return
	}

	type res struct {
		Message string `json:"message"`
	}
	response:=res {
		Message: fmt.Sprintf("%v has been updated in the record", body.Url),
	}
	sendAsJson(w, response)
}

func (tCtx TinyCtx) tinyGetHandler(w http.ResponseWriter, r *http.Request) {
	url:=r.URL.Query().Get("url")
	resultList, err := tCtx.client.HGetAll(tCtx.ctx, "urlHash").Result()
	if err != nil {
		redisErr(err)
		return
	}
	type res struct {
		Result map[string]string `json:"result"`
	}

	for val, key := range(resultList) {
		if val == url {
			response:=res {
				map[string]string{url: key },
			}
			sendAsJson(w, response)
		return	
		}
	}
	sendErrJson(w, fmt.Sprintf("%v not found in record", url))	
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

func (tCtx TinyCtx) tinyDelHandlerByHash(w http.ResponseWriter, r *http.Request) {
	hash:=r.URL.Query().Get("hash")
	result, err :=tCtx.client.HDel(tCtx.ctx, "urlHash", hash).Result()
	fmt.Println("res", result)
	if err != nil {
		redisErr(err)
		return
	}
	if result == 0 {
		sendErrJson(w, fmt.Sprintf("%v hash does not exist!", hash))
		return 
	}

	type res struct {
		Message string `json:"message"`
	}
	response:=res {
		Message: fmt.Sprintf("%v has been deleted from the record", hash),
	}
	sendAsJson(w, response)
}


func main(){
	fmt.Println("welcome to the tiny app")
	r:= mux.NewRouter()	
	client := redis.NewClient(&redis.Options{

        Addr:	  "redis:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })
	ctx := context.Background()

	tCtx :=TinyCtx{
		client: client,
		ctx: ctx,
	}
	r.HandleFunc("/{hash}",tCtx.rootRedirector).Methods("GET")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyPostHandler).Methods("POST")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyPutHandler).Methods("PUT")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyGetHandler).Methods("GET")

	r.HandleFunc("/api/v2/tiny/all", tCtx.tinyGetAllHandler).Methods("GET")

	r.HandleFunc("/api/v1/tiny", tCtx.tinyDelHandlerByHash).Methods("DELETE")
	
	cli.CliApp();
	http.ListenAndServe(":6969",r)
}
