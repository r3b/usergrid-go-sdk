package usergrid

import (
	// "github.com/user/newmath"
	"testing"
	"encoding/json"
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
