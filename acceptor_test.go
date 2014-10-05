package main

import(  
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/julienschmidt/httprouter"
)

// Test basic responder flow that returns hello world
func Test_GetHelloWorld(t *testing.T) {
   
  router := httprouter.New()
  router.GET("/_/*p", process)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/_/test_string", nil)    

  _start_broker()
  go _test_string_responder()
  
  router.ServeHTTP(writer, request)
  
  if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}  
  if writer.Body.String() != "Hello World" {
    t.Errorf("Body is %v", writer.Body)
  }  
}

// Test basic responder flow that returns a JSON string that it sends across
func Test_GetJSON(t *testing.T) {
  router := httprouter.New()
  router.GET("/_/*p", process)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/_/test_json", nil)    

  _start_broker()
  go _test_json_responder()
  
  router.ServeHTTP(writer, request)
  
  if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}  
  sample_json := _read_sample_json()
  if writer.Body.String() != sample_json {
    t.Errorf("Body is %s", writer.Body)
  }    
}

