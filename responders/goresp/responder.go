package main

import (
	zmq "github.com/pebbe/zmq4"
	"fmt"
  "code.google.com/p/go-uuid/uuid"
)

const (
	ROUTEID = "GET/_/hello/go"
)

func main() {
	responder, _ := zmq.NewSocket(zmq.REQ)
	defer responder.Close()

	identity := uuid.New()
	responder.SetIdentity(identity)
	responder.Connect("tcp://localhost:4321")

	//  Tell broker we're ready for work
	fmt.Printf("%s - %s responder ready\n", ROUTEID, identity)
	responder.Send(ROUTEID, 0)

	for {
		msg, err := responder.RecvMessage(0)
		if err != nil {
      fmt.Println("Error in receiving message:", err)
			break //  Interrupted
		}
    resp := []string{"200", "{\"Content-Type\": \"text/html\"}", msg[0],}
    fmt.Println("Responding with:", resp)
		responder.SendMessage(ROUTEID, resp)
	}
}
