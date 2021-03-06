package main

import (
    "flag"
    "log"
    "fmt"
    "time"
    "os"
    "strconv"
    "encoding/json"
    "github.com/r3b/usergrid-go-sdk"
)

var (
    API string
    PAGE_SIZE int
    ORGNAME string
    APPNAME string
    CLIENT_ID string
    CLIENT_SECRET string
    ENDPOINT string
    REQUESTS int
    RESPONSES int
    RESPONSE_SIZE int
)

func init(){
	flag.StringVar(&API, "apiurl", "http://api.usergrid.com", "Usergrid base API URI")
	flag.StringVar(&ORGNAME, "orgname", "yourorgname", "The name of your org")
	flag.StringVar(&APPNAME, "appname", "sandbox", "The name of your app")
	flag.StringVar(&CLIENT_ID, "id", "", "The Client ID for your app")
	flag.StringVar(&CLIENT_SECRET, "secret", "", "The Client Secret for your app")
	flag.StringVar(&ENDPOINT, "endpoint", "", "The endpoint to fetch")
	flag.IntVar(&PAGE_SIZE, "pagesize", 10, "The number of items in history to be displayed")
	flag.Parse()
	log.SetOutput(os.Stderr)
}

func FetchAll(client usergrid.Client, endpoint string, entities chan<- interface{}, control chan<- bool, cursor string) {
	var objmap interface{}
	params := map[string]string{"limit":strconv.Itoa(PAGE_SIZE)}
	if cursor!="" {
		params["cursor"]=cursor
	}
	err:= client.Get(endpoint,params, usergrid.JSONResponseHandler(&objmap))
	if err != nil {
		log.Printf("ERROR: %s\n\n", err)
		control <- true
		// return
	}else{
		obmap:=objmap.(map[string]interface{})
		if len(obmap["entities"].([]interface{}))>0 {
			for _,v := range obmap["entities"].([]interface{}) {
				entities <- v
			}
		}
		if(obmap["cursor"]!=nil){
			cursor,_ := obmap["cursor"].(string)
			go FetchAll(client, endpoint, entities, control, cursor)
		}else{
			control <- true
		}
	}

}
func main(){
	client := usergrid.Client {Organization:ORGNAME,Application:APPNAME,Uri:API}
	var objmap interface{}
	var results = make([]interface{},0)
	messages := make(chan interface{})
	done := make(chan bool, 1)
	if(len(CLIENT_ID)>0 && len(CLIENT_SECRET)>0){
		err := client.OrgLogin(CLIENT_ID, CLIENT_SECRET)
		if err != nil {
			log.Printf(err.Error());
			return
		}
	}
	go func(){
		for {
	        select {
	        case v := <-messages:
	        	results=append(results, v)
	        case <-time.After(time.Second * 10):
	        	fmt.Println("Timeout!")
	        }
	    }
	}()
	client.Get("",nil, usergrid.JSONResponseHandler(&objmap))
	str,_:=json.MarshalIndent(objmap,"","  ")
	log.Printf("RESPONSE: %s", str)
	go FetchAll(client, ENDPOINT, messages, done, "")
	<- done
	entities, _ := json.MarshalIndent(results, "", "  ")
	fmt.Printf("%s\n", entities)
	log.Printf("Done. %d requests and %d responses", REQUESTS, RESPONSES)
    return
}

