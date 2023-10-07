package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func callApi[K comparable, V any](method string, url string, body map[K]V)  {
    jBody, err := json.Marshal(body)
    req, err := http.NewRequest(method, url, bytes.NewBuffer(jBody)) 
    if err != nil {
	log.Panic("Failed to get response")
    } 
    client := &http.Client{}
    resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

    defer resp.Body.Close()
    print(resp.Body)
    // body, err =io.ReadAll(resp.Body)
 //    if err != nil {
	// log.Panic("can not able to parse body")
 //    }
 //    var resObj ResResult
 //    json.Unmarshal(body, &resObj)
 //    return resObj
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
		    fmt.Println("this is the url", ctx.Args().First())
		    return nil
		},
	    },
	    {
		Name: "update",
		Aliases: []string{"u"},
		Usage: "update the url",
		Action: func(ctx *pkcli.Context) error {
		    fmt.Println("this is the url", ctx.Args().First())
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
		    res:=callApi("GET", "http://localhost:6969/api/v2/tiny/all", body)
		    w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.Debug)
		    fmt.Fprintf(w, "%v\t%v\n\n", "short url", "original url")
		    for key, val :=range(res.Result) {
			fmt.Fprintf(w, "%v\t%v\n", key, val)
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
