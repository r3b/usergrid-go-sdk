package usergrid

import (
	// "github.com/user/newmath"
	"testing"
	"encoding/json"
	"strconv"
    "fmt"
    // "time"
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
func TestMassDelete(t *testing.T){
	client := Client {Organization:"yourorgname",Application:"sandbox",Uri:"https://api.usergrid.com"}
	params := map[string]string{"limit":strconv.Itoa(500)}
	resp, err:= client.Get("benchmark", params)
	if(err!=nil){
		t.Logf("Test failed: %s\n", err)
		t.Fail()
	}
	// if(t.Failed()){
	// 	str,_:=json.MarshalIndent(resp,"","  ")
	// 	t.Logf("RESPONSE: %s", str)
	// }
	if (resp["entities"]!=nil && len(resp["entities"].([]interface{}))>0) {
		for k,v := range resp["entities"].([]interface{}) {
			if(v == nil){
				t.Logf("could not delete %d: %v\n", k, v)
				t.Fail()
			}else{
				t.Logf(fmt.Sprintf("yay %v",v))
				entity := v.(map[string]interface{})
				t.Logf(fmt.Sprintf("entity %v",entity["uuid"]))
				_, err:= client.Delete("benchmark/"+entity["uuid"].(string), nil)
				if(err!=nil){
					t.Logf("could not delete %s: %s\n", entity["uuid"].(string), err)
					t.Fail()
				}
			}
		}
	}
	// resp:= client.RequestChannel("GET", "http://api.usergrid.com/yourorgname/sandbox/benchmarks",nil, nil)
	
 //    select {
 //    case v := <-resp:
 //    	str,_:=json.MarshalIndent(v,"","  ")
 //    	fmt.Printf("%s\n", str)
 //    case <-time.After(time.Second * 20):
 //    	panic("Timeout!")
 //    }
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

