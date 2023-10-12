package main

import (
	"bytes"
	"encoding/json"
	"fmt"
    /* "io" */
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	pkcli "github.com/urfave/cli/v2"
)
type ResResult struct {
    Result map[string]string `json:"result"`
}

// type ApiModel struct {
//     url string
//     method string
//     body map[]
// }

// type httpMethods interface {
//     http.MethodPost | http.MethodGet
// }

type CreateUrl struct {
    Url string `json:"url"`
    Expitry int `json:"expiry"`
}

type UpdateUrl struct {
    Url string `json:"url"`
    OldUrl string `json:"oldurl"`
    Expitry int `json:"expiry"`
}

type ResMessage struct {
    Message string `json:"message"`
}

func callApi(method string, url string, body interface{}, decoder interface{}) interface{} {
    jBody, err := json.Marshal(body)
    req, err := http.NewRequest(method, url, bytes.NewBuffer(jBody)) 
    if err != nil {
	log.Panic("Failed to get response")
    } 
    res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

    defer res.Body.Close()
 //    resBody, err :=io.ReadAll(res.Body)
 //    if err != nil {
	// log.Panic("can not able to parse body")
 //    }
    return json.NewDecoder(res.Body).Decode(decoder)
}

func CliApp() {
    app:= &pkcli.App{
	Name: "tiny",
	Usage: "make tiny urls with just a click",
	Action: func(*pkcli.Context) error{
	    fmt.Println("this is tiny cli app")
	    return nil
	},
	Commands: []*pkcli.Command{
	    {
		Name: "make",
		Aliases: []string{"m"},
		Usage: "Make a tiny url",
		Action: func(ctx *pkcli.Context) error {
		    body:= CreateUrl{
			Url: ctx.Args().First(),
			Expitry: 5,	
		    }
		    var res ResMessage 
		    callApi("POST", "http://localhost:6969/api/v1/tiny", body, &res)
		    fmt.Println(res.Message)
		    return nil
		},
	    },
	    {
		Name: "update",
		Aliases: []string{"u"},
		Usage: "update the url",
		Action: func(ctx *pkcli.Context) error {
		    body:= UpdateUrl{
			Url: ctx.Args().First(),
			OldUrl: ctx.Args().Tail()[0],
			Expitry: 5,	
		    }
		    var res ResMessage 
		    callApi("PUT", "http://localhost:6969/api/v1/tiny", body, &res)
		    fmt.Println(res.Message)
		    return nil
		},
	    },
	    {
		Name: "list",
		Aliases: []string{"l"},
		Usage: "list all tiny urls",
		Action: func(ctx *pkcli.Context) error {
		    fmt.Printf("these are the list of urls\n\n")
		    var body map[string]string
		    var res ResResult
		    err := callApi("GET", "http://localhost:6969/api/v1/tiny/all", body, &res)
		    if err != nil {
			fmt.Println(err)
		    }
		    w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.Debug)
		    fmt.Fprintf(w, "%v\t %v\n\n", "short url", "original url")
		    for key, val :=range(res.Result) {
			fmt.Fprintf(w, "%v\t %v\n", key, val)
		    }
		    w.Flush()
		    return nil
		},
	    },
	},
    }
    if err := app.Run(os.Args); err != nil {
	log.Fatal(err)
    }
}

func main(){
    CliApp()
}
