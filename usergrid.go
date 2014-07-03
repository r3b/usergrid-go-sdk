package usergrid

import (
    // "flag"
    "log"
    // "errors"
    "net/http"
    "net/url"
	"strings"
    "io"
    "io/ioutil"
    //"text/template"
    "fmt"
    "time"
    // "strings"
    "encoding/json"
    // "code.google.com/p/go-uuid/uuid"
    // "runtime"                                                                                                                              
    "strconv"
    "os"
    // "os/signal"
    // "syscall"
    "reflect"
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
	// flag.StringVar(&API, "apiurl", "http://api.usergrid.com", "Usergrid base API URI")
	// flag.StringVar(&ORGNAME, "orgname", "yourorgname", "The name of your org")
	// flag.StringVar(&APPNAME, "appname", "sandbox", "The name of your app")
	// flag.StringVar(&CLIENT_ID, "id", "", "The Client ID for your app")
	// flag.StringVar(&CLIENT_SECRET, "secret", "", "The Client Secret for your app")
	// flag.StringVar(&ENDPOINT, "endpoint", "", "The endpoint to fetch")
	// flag.IntVar(&PAGE_SIZE, "pagesize", 10, "The number of items in history to be displayed")
	// flag.Parse()
	log.SetOutput(os.Stderr)
}
type Client struct {
	Organization,Application,Uri,access_token string
	//,grant_type,client_id,client_secret
}
type ResponseHandlerInterface func(body io.ReadCloser) error
func JSONResponseHandler(objmap *interface{}) (ResponseHandlerInterface){
	return func(body io.ReadCloser) (error){
		responseBody, _ := ioutil.ReadAll(body)
		if err := json.Unmarshal(responseBody, &objmap); err == nil{
			err:=CheckForError(*objmap)
			return err
		}else{
			return err
		}
		return nil
	}
}
func CheckForError(objmap interface{}) (error){
	omap:=objmap.(map[string]interface{})
	str := ""
	if omap["error"] != nil {
		if omap["error_description"] != nil {
			str = omap["error_description"].(string)
		}else{
			str = omap["error"].(string)
		}
		// log.Printf("an error was returned: %v\n", str)
		return &UsergridError{Message:str}
	}
	return nil	
}

func PrintAll(vals []interface{}) {
    for k, v := range vals {
        fmt.Println(k, reflect.TypeOf(v), v)
    }
}
func AppendQueryParams(endpoint string, params map[string]string) string{
    u, _ := url.Parse(endpoint)
    if params != nil {
    	q := u.Query()
	    for k, v := range params {
	    	q.Set(k,v)
		}
		u.RawQuery=q.Encode()
		endpoint = fmt.Sprintf("%s?%s",endpoint,u.RawQuery)
    }
    return endpoint
}
func (client *Client) Authenticate(grant_type string, client_id string, client_secret string){}
func (client *Client) Login(username string, password string) error{
	urlStr := fmt.Sprintf("%s/%s/%s/%s/%s/%s",client.Uri,client.Organization, client.Application, "users", username,  "token")
	data := map[string]string{"grant_type":"password","username":username,"password":password}
	rawmap, err := client.Request("POST", urlStr, nil, data)
	if err == nil {
		client.access_token, err = RawToString(rawmap["access_token"])
	}
	return err
}
func (client *Client) OrgLogin(client_id string, client_secret string) error {
	urlStr := fmt.Sprintf("%s/%s",client.Uri,"management/token")
	data := map[string]string{"grant_type":"client_credentials","client_id":client_id,"client_secret":client_secret}
	var objmap interface{}
	err := client.RequestAsync("POST", urlStr, nil, data, func(body io.ReadCloser) (error){
		responseBody, _ := ioutil.ReadAll(body)
		err := json.Unmarshal(responseBody, &objmap)
		omap:=objmap.(map[string]interface{})
		client.access_token=omap["access_token"].(string)
		err=CheckForError(objmap)
		return err

	})
	return err
}
func (client *Client) AddAuthorizationHeaders(req *http.Request){
	if(len(client.access_token) > 0){
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",client.access_token))	
	}
}
func (client *Client) MakeRequest(method string, endpoint string, params map[string]string, data interface{}) (*http.Request, error) {
	var err error
	var req *http.Request
	method = strings.ToUpper(method)
    endpoint = AppendQueryParams(endpoint, params)
    switch strings.ToUpper(method) {
	    case "POST":
		    body, _ := json.Marshal(data)
		    req, err = http.NewRequest(method, endpoint, strings.NewReader(string(body)))
		case "PUT":
		    body, _ := json.Marshal(data)
		    req, err = http.NewRequest(method, endpoint, strings.NewReader(string(body)))
	    case "DELETE":
	    	req, err = http.NewRequest(method, endpoint, nil)
	    default: //GET
	    	method="GET"
	    	req, err = http.NewRequest(method, endpoint, nil)
    }
    if err != nil {
    	return nil, err
    }
    client.AddAuthorizationHeaders(req)
    return req, nil
}


func (client *Client) RequestAsync(method string, endpoint string, params map[string]string, data interface{}, handler ResponseHandlerInterface) (error) {
	_client := &http.Client{}
	req, err :=client.MakeRequest(method, endpoint, params, data)
	if err != nil {
	    return err
	}
    REQUESTS++
    resp, err :=_client.Do(req)
	defer resp.Body.Close()
	if err != nil {
	    return err
	}
	RESPONSES++
	return handler(resp.Body)
}
func (client *Client) Request(method string, endpoint string, params map[string]string, data interface{}) (map[string]*json.RawMessage, error) {
	_client := &http.Client{}
	req, err :=client.MakeRequest(method, endpoint, params, data)
	if err != nil {
	    // log.Panicf("ERROR %s\n", err)
	    return nil, err
	}
    REQUESTS++
    resp, err :=_client.Do(req)
	defer resp.Body.Close()
	if err != nil {
	    // log.Panicf("ERROR %s\n", err)
	    return nil, err
	}
	RESPONSES++
	// log.Printf("%d\t%s\t%s\n", resp.StatusCode, method, endpoint)
	// responseBody, err := ioutil.ReadAll(resp.Body)
	var objmap map[string]*json.RawMessage
	decoder:=json.NewDecoder(resp.Body)
	err=decoder.Decode(&objmap)
	if err != nil {
		return nil, err
	}else if err:=CheckForError(objmap); err != nil {
		return nil, err
	}
	// if err != nil {
	// 	return nil, err
	// }
	// err = client.CheckForError(responseBody)
	return objmap, nil
}
func (client *Client) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestAsync("GET",urlStr, params, nil, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
func (client *Client) Delete(endpoint string, params map[string]string) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestAsync("DELETE",urlStr, params, nil, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
func (client *Client) Post(endpoint string, params map[string]string, data interface{}) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestAsync("POST",urlStr, params, data, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
func (client *Client) Put(endpoint string, params map[string]string, data interface{}) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestAsync("PUT",urlStr, params, data, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
type UsergridError struct {
	error, timestamp, duration, exception, Message string
}
func (err *UsergridError) Error() string {
	return err.Message
}
func RawToString(raw *json.RawMessage) (string,error) {
	var str string
	if *raw == nil {
		return "", nil
	}
	err := json.Unmarshal(*raw, &str)
	return str, err
}


func (client *Client) GetAll(endpoint string, entities chan<- interface{}, control chan<- bool, cursor string) {
	params := map[string]string{"limit":strconv.Itoa(PAGE_SIZE)}
	if cursor!="" {
		params["cursor"]=cursor
	}
	resp,err:= client.Get(endpoint,params)
	if err != nil {
		log.Printf("ERROR: %s\n\n", err)
		control <- true
		// return
	}else{
		if(resp["cursor"]!=nil){
			cursor,_ := resp["cursor"].(string)
			go client.GetAll(endpoint, entities, control, cursor)
		}else{
			control <- true
		}
		if len(resp["entities"].([]interface{}))>0 {
			for _,v := range resp["entities"].([]interface{}) {
				entities <- v
			}
		}
	}

}
func main(){
	client := Client {Organization:ORGNAME,Application:APPNAME,Uri:API}
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
	        // default:
	        	// fmt.Println("Waiting...")
	            // fmt.Println("received", v)
	        }
	    }
	}()
	go client.GetAll(ENDPOINT, messages, done, "")
	<- done
	entities, _ := json.MarshalIndent(results, "", "  ")
	fmt.Printf("%s\n", entities)
	log.Printf("Done. %d requests and %d responses", REQUESTS, RESPONSES)

    // log.Printf("Done. %d requests and %d responses", REQUESTS, RESPONSES)
    return
}

