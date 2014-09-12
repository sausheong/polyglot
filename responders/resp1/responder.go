package main
 
import (
    zmq "github.com/pebbe/zmq4"
)
 
func main() {
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.SetIdentity("1123581321") //key to be used
	responder.Connect("tcp://127.0.0.1:10001")
 
	for {
		msg, _ := responder.RecvMessage(0)
		reqId := msg[0]
		data := msg[1]
		// data contains "Hello, World"
 
		result := []string{reqId, string("Hello, ZMQ")}
		responder.SendMessage(result, 0)
	}
}