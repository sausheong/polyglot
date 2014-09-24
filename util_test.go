package main

import(  
  	zmq "github.com/pebbe/zmq4"
  	"fmt"
    "code.google.com/p/go-uuid/uuid"    
    "os/exec"
    "io/ioutil"
)

// 
// Responders
// 

func _test_string_responder() {
  routeid := "GET/_/test_string"
	responder, _ := zmq.NewSocket(zmq.REQ)
	defer responder.Close()

	identity := uuid.New()
	responder.SetIdentity(identity)
	responder.Connect("tcp://localhost:4321")
	responder.Send(routeid, 0)

	for {
		_, err := responder.RecvMessage(0)
		if err != nil {
      fmt.Println("Error in receiving message:", err)
			break //  Interrupted
		}
    resp := []string{"200", "{\"Content-Type\": \"text/html\"}", "Hello World",}
		responder.SendMessage(routeid, resp)
	}
}

func _test_json_responder() {
  routeid := "GET/_/test_json"
	responder, _ := zmq.NewSocket(zmq.REQ)
	defer responder.Close()

	identity := uuid.New()
	responder.SetIdentity(identity)
	responder.Connect("tcp://localhost:4321")
	responder.Send(routeid, 0)

  json := _read_sample_json()
	for {
		_, err := responder.RecvMessage(0)
		if err != nil {
      fmt.Println("Error in receiving message:", err)
			break //  Interrupted
		}
    resp := []string{"200", "{\"Content-Type\": \"application/json\"}", json,}
		responder.SendMessage(routeid, resp)
	}
}

func _start_broker() {
  broker := exec.Command("go run ./broker/broker.go") 
  broker.Start()
}

func _read_sample_json() string {
  json_bytes, err := ioutil.ReadFile("./testdata/sample.json")
  if err != nil {
    fmt.Println("Error in reading sample json", err)
  }
  return string(json_bytes)  
}