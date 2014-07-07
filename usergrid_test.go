package usergrid

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)
func TestPost(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	data := map[string]string{"name":"go_dog", "status":"good dog"}
	resp, err:= client.Post("dogs", nil, data)
	if(err!=nil){
		t.Logf("TestPost failed: %s\n", err)
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestPostDuplicate(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	data := map[string]string{"name":"go_dog", "status":"good dog"}
	resp, err:= client.Post("dogs", nil, data)
	if(err==nil){
		t.Logf("TestPost failed. Should have received a duplicate entity error.\n")
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestGet(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	resp, err:= client.Get("dogs/go_dog", nil)
	if(err!=nil){
		t.Logf("TestGet failed: %s\n", err)
		t.Fail()
	}else{
		if (resp["entities"]==nil) {
			t.Logf("No entities returned")
			t.Fail()
		}else if(len(resp["entities"].([]interface{}))!=1){
			t.Logf("Incorrect number of entities")
			t.Fail()
		}
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestGetBadEntity(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	resp, err:= client.Get("dogs/go_dogizzle", nil)
	if(err==nil){
		t.Logf("TestGet failed. dog should not exist: %s\n", err)
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestPut(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	data := map[string]string{"status":"flat out on the highway and drying in the sun"}
	resp, err:= client.Put("dogs/go_dog", nil, data)
	if(err!=nil){
		t.Logf("TestPut failed: %s\n", err)
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestDelete(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	resp, err:= client.Delete("dogs/go_dog", nil)
	if(err!=nil){
		t.Logf("Test failed: %s\n", err)
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestDelete2(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	resp, err:= client.Get("dogs/go_dog", nil)
	if(err==nil){
		t.Logf("Test failed. Dog should be deleted: %s\n", err)
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}
func TestDelete3(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	resp, err:= client.Delete("dogs/go_dog", nil)
	if(err==nil){
		t.Logf("Test failed. Should return service_resource_not_found error\n")
		t.Fail()
	}
	if(t.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		t.Logf("RESPONSE: %s", str)
	}
}

func BenchmarkPost(b *testing.B) {
 	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	data := map[string]string{"index":"0", "description":"golang benchmark"}
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
    	data["index"]=strconv.Itoa(i);
		resp, err:= client.Post("benchmark", nil, data)
		if(err!=nil){
			b.Logf("BenchmarkPost failed: %s\n", err)
			b.Fail()
		}
		if(b.Failed()){
			str,_:=json.MarshalIndent(resp,"","  ")
			b.Logf("RESPONSE: %s", str)
		}
   }
}
func BenchmarkRawRequests(b *testing.B) {
 	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	params := map[string]string{"limit":strconv.Itoa(10)}
	resp, err:= client.Get("benchmark", params)
	if(err!=nil){
		b.Logf("Test failed: %s\n", err)
		b.Fail()
	}
	if (resp["entities"]==nil || len(resp["entities"].([]interface{}))==0) {
		b.Logf("Test failed: no entities\n")
		b.Fail()
	}
	entities:=resp["entities"].([]interface{})
	requests:=0
	responses:=0
	errors:=0
	var objmap interface{}
	var entity map[string]interface{}
	responseChan := make(chan []byte)
	go func(){
		for {
	        select {
	        case v := <-responseChan:
	        	if err := json.Unmarshal(v, &objmap); err == nil{
					if err:=CheckForError(&objmap); err != nil {
						errors++
					}
				}else{
					errors++
				}
	        	
	        	responses++
	        	requests--
	        case <-time.After(time.Second * 10):
	        	return
	        }
	    }
	}()
 	b.ResetTimer()
    for i := 0; i < b.N; i++ {
    	if(len(entities)==0){
			b.Logf("Test failed: we ran out of entities\n")
			b.Fail()
			continue
    	}
    	entity = entities[i % len(entities)].(map[string]interface{})
    	for ;requests>=MAX_CONCURRENT_REQUESTS;{
        	// if we outpace GOMAXPROCS, we'll run out of threads
    		time.Sleep(60 * time.Millisecond)
    	}
		client.Request("GET", "http://api.usergrid.com/yourorgname/sandbox/benchmark/"+entity["uuid"].(string), nil, nil, responseChan)
		requests++
   }
   for ;requests>0;{
		// wait for the last few responses
		time.Sleep(60 * time.Millisecond)   	
   }
}
func BenchmarkGet(b *testing.B) {
 	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	params := map[string]string{"limit":strconv.Itoa(500)}
	resp, err:= client.Get("benchmark", params)
	if(err!=nil){
		b.Logf("Test failed: %s\n", err)
		b.Fail()
	}
	if (resp["entities"]==nil || len(resp["entities"].([]interface{}))==0) {
		b.Logf("Test failed: no entities to delete\n")
		b.Fail()
	}
	entities:=resp["entities"].([]interface{})
	var entity map[string]interface{}
 	b.ResetTimer()
    for i := 0; i < b.N; i++ {
    	if(len(entities)==0){
			b.Logf("Test failed: we ran out of entities\n")
			b.Fail()
			continue
    	}
    	entity, entities = entities[len(entities)-1].(map[string]interface{}), entities[:len(entities)-1]
		_, err:= client.Get("benchmark/"+entity["uuid"].(string), nil)
		if(err!=nil){
			b.Logf("BenchmarkGet failed: %s\n", err)
			b.Fail()
		}
   }
	if(b.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		b.Logf("RESPONSE: %s", str)
	}
}
func BenchmarkPut(b *testing.B) {
 	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	params := map[string]string{"limit":strconv.Itoa(500)}
	resp, err:= client.Get("benchmark", params)
	if(err!=nil){
		b.Logf("Test failed: %s\n", err)
		b.Fail()
	}
	if (resp["entities"]==nil || len(resp["entities"].([]interface{}))==0) {
		b.Logf("Test failed: no entities to delete\n")
		b.Fail()
	}
	entities:=resp["entities"].([]interface{})
	var entity map[string]interface{}
 	b.ResetTimer()
    for i := 0; i < b.N; i++ {
    	if(len(entities)==0){
			b.Logf("Test failed: we ran out of entities\n")
			b.Fail()
			continue
    	}
    	entity, entities = entities[len(entities)-1].(map[string]interface{}), entities[:len(entities)-1]
		_, err:= client.Put("benchmark/"+entity["uuid"].(string), nil, map[string]string{"updated":"true"})
		if(err!=nil){
			b.Logf("BenchmarkPut failed: %s\n", err)
			b.Fail()
		}
   }
	if(b.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		b.Logf("RESPONSE: %s", str)
	}
}

func BenchmarkDelete(b *testing.B) {
 	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	params := map[string]string{"limit":strconv.Itoa(500)}
	resp, err:= client.Get("benchmark", params)
	if(err!=nil){
		b.Logf("Test failed: %s\n", err)
		b.Fail()
	}
	if (resp["entities"]==nil || len(resp["entities"].([]interface{}))==0) {
		b.Logf("Test failed: no entities to delete\n")
		b.Fail()
	}
	entities:=resp["entities"].([]interface{})
	var entity map[string]interface{}
 	b.ResetTimer()
    for i := 0; i < b.N; i++ {
    	if(len(entities)==0){
			b.Logf("Test failed: we ran out of entities\n")
			b.Fail()
			continue
    	}
    	entity, entities = entities[len(entities)-1].(map[string]interface{}), entities[:len(entities)-1]
		_, err:= client.Delete("benchmark/"+entity["uuid"].(string), nil)
		if(err!=nil){
			b.Logf("BenchmarkDelete failed: %s\n", err)
			b.Fail()
		}
   }
	if(b.Failed()){
		str,_:=json.MarshalIndent(resp,"","  ")
		b.Logf("RESPONSE: %s", str)
	}
}

