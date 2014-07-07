package usergrid

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "reflect"
    "runtime"
	"strings"
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
    MAX_CONCURRENT_REQUESTS int = runtime.NumCPU()
)

func init(){
	log.SetOutput(os.Stderr)
	runtime.GOMAXPROCS(MAX_CONCURRENT_REQUESTS)
}
type Client struct {
	Organization,Application,Uri,access_token string
	_client *http.Client
}
type ResponseHandlerInterface func(responseBody []byte) error
func JSONResponseHandler(objmap *interface{}) (ResponseHandlerInterface){
	return func(responseBody []byte) (error){
		if err := json.Unmarshal(responseBody, &objmap); err == nil{
			err:=CheckForError(objmap)
			return err
		}else{
			return err
		}
		return nil
	}
}
func CheckForError(objmap *interface{}) (error){

	omap:=(*objmap).(map[string]interface{})
	str := ""
	if omap["error"] != nil {
		if omap["error_description"] != nil {
			str = omap["error_description"].(string)
		}else{
			str = omap["error"].(string)
		}
		return errors.New(str)
	}
	return nil	
}

func PrintAll(vals []interface{}) {
    for k, v := range vals {
        log.Println(k, reflect.TypeOf(v), v)
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
	var objmap interface{}
	err := client.RequestWithHandler("POST", urlStr, nil, data, func(responseBody []byte) (error){
		err := json.Unmarshal(responseBody, &objmap)
		omap:=objmap.(map[string]interface{})
		client.access_token=omap["access_token"].(string)
		err=CheckForError(&objmap)
		return err

	})
	return err
}
func (client *Client) OrgLogin(client_id string, client_secret string) error {
	urlStr := fmt.Sprintf("%s/%s",client.Uri,"management/token")
	data := map[string]string{"grant_type":"client_credentials","client_id":client_id,"client_secret":client_secret}
	var objmap interface{}
	err := client.RequestWithHandler("POST", urlStr, nil, data, func(responseBody []byte) (error){
		err := json.Unmarshal(responseBody, &objmap)
		omap:=objmap.(map[string]interface{})
		client.access_token=omap["access_token"].(string)
		err=CheckForError(&objmap)
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


func (client *Client) Request(method string, endpoint string, params map[string]string, data interface{}, responseChan chan []byte){
	//intialize an http client if we don't already have one
	if client._client == nil {
		client._client= &http.Client{}
	}
	defer func() {
        if r := recover(); r != nil {
            responseChan <- []byte(fmt.Sprintf("{\"error\":\"%s\", \"error_description\":\"%s: %v\"}", "network_error", "The request failed at the network level", r))
        }
    }()
	go func(){
		req, err :=client.MakeRequest(method, endpoint, params, data)
		if err != nil {
		    log.Panic(err)
		}
	    REQUESTS++
	    resp, err :=client._client.Do(req)
	    if err != nil {
		    log.Panic(err)
		}
		defer resp.Body.Close()
		RESPONSES++
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}
		responseChan <- responseBody	
	}()
}
func (client *Client) RequestWithHandler(method string, endpoint string, params map[string]string, data interface{}, handler ResponseHandlerInterface) (error) {
    responseChan := make(chan []byte)
	client.Request(method, endpoint, params, data, responseChan)
	responseBody := <-responseChan
	return handler(responseBody)
}

func (client *Client) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestWithHandler("GET",urlStr, params, nil, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
func (client *Client) Delete(endpoint string, params map[string]string) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestWithHandler("DELETE",urlStr, params, nil, JSONResponseHandler(&objmap))
	if(objmap == nil){
		return nil, err
	}else{
		return objmap.(map[string]interface{}), err
	}
}
func (client *Client) Post(endpoint string, params map[string]string, data interface{}) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestWithHandler("POST",urlStr, params, data, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
func (client *Client) Put(endpoint string, params map[string]string, data interface{}) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("%s/%s/%s/%s",client.Uri,client.Organization, client.Application, endpoint);
	var objmap interface{}
	err := client.RequestWithHandler("PUT",urlStr, params, data, JSONResponseHandler(&objmap))
	return objmap.(map[string]interface{}), err
}
