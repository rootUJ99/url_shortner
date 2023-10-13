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

type CreateUrl struct {
    Url string `json:"url"`
    Expitry int `json:"expiry"`
}
type DeleteUrl struct {
    Url string `json:"url"`
}

type UpdateUrl struct {
    Url string `json:"url"`
    OldUrl string `json:"oldurl"`
    Expitry int `json:"expiry"`
}

type ResMessage struct {
    Message string `json:"message"`
}

const host string = "http://localhost:6969"

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
		Name: "create",
		Aliases: []string{"c"},
		Usage: "create a new tiny url",
		Action: func(ctx *pkcli.Context) error {
		    body:= CreateUrl{
			Url: ctx.Args().First(),
			Expitry: 5,	
		    }
		    var res ResMessage 
		    endPoint:= fmt.Sprintf("%v/api/v1/tiny", host)
		    callApi("POST", endPoint, body, &res)
		    fmt.Println(res.Message)
		    return nil
		},
	    },
	    {
		Name: "update",
		Aliases: []string{"u"},
		Usage: "update the existing tiny url",
		Action: func(ctx *pkcli.Context) error {
		    body:= UpdateUrl{
			OldUrl: ctx.Args().First(),
			Url: ctx.Args().Tail()[0],
			Expitry: 5,	
		    }
		    fmt.Println(body)
		    var res ResMessage 
		    endPoint:= fmt.Sprintf("%v/api/v1/tiny", host)
		    callApi("PUT", endPoint, body, &res)
		    fmt.Println(res.Message)
		    return nil
		},
	    },
	    {
		Name: "delete",
		Aliases: []string{"d"},
		Usage: "delete an existing tiny url",
		Action: func(ctx *pkcli.Context) error {
		    body:= DeleteUrl{
			Url: ctx.Args().First(),
		    }
		    var res ResMessage 
		    endPoint:= fmt.Sprintf("%v/api/v1/tiny", host)
		    callApi("DELETE", endPoint, body, &res)
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
		    endPoint:= fmt.Sprintf("%v/api/v1/tiny/all", host)
		    err := callApi("GET", endPoint, body, &res)
		    if err != nil {
			fmt.Println(err)
		    }
		    w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.Debug)
		    fmt.Fprintf(w, "%v\t %v\n\n", "short url", "original url")
		    for key, val :=range(res.Result) {
			urlWithHost := fmt.Sprintf("%v/%v", host, key)
			fmt.Fprintf(w, "%v\t %v\n", urlWithHost, val)
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
