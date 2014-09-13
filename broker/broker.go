package main

import (
  zmq "github.com/pebbe/zmq4"
  "fmt"
  "math/rand"
  "time"
  "log"
  "os"
)

var logger *log.Logger
var responders map[string][]string
var parked map[string][]string

func init() {
  file, err := os.OpenFile("broker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
      log.Fatalln("Failed to open log file", err)
  }  
  logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
  fmt.Println("Polyglot Broker v0.1")
  fmt.Println("")
  frontend, _ := zmq.NewSocket(zmq.ROUTER)
  backend, _  := zmq.NewSocket(zmq.ROUTER)
  admin, _    := zmq.NewSocket(zmq.ROUTER)
  defer frontend.Close()
  defer backend.Close()
  defer admin.Close()
  frontend.Bind("tcp://*:1234") //  For clients
  backend.Bind("tcp://*:4321")  //  For responders
  admin.Bind("tcp://*:9999")    //  For admin

  //  Queue of available responders
  responders = make(map[string][]string)
  
  // parked responders
  parked = make(map[string][]string)
  
  poller := zmq.NewPoller()
  poller.Add(backend, zmq.POLLIN)
  poller.Add(frontend, zmq.POLLIN)
  poller.Add(admin, zmq.POLLIN)

LOOP:
  for {
    //  Poll frontend only if we have available responders
    var sockets []zmq.Polled
    var err error
    sockets, err = poller.Poll(-1)
    if err != nil {
      break //  Interrupted
    }
    for _, socket := range sockets {
      
      switch s := socket.Socket; s {
      
      // Communicate with admin
      case admin:
        msg, err := s.RecvMessage(0)
        if err != nil {
          danger("Admin", "Cannot receive message", err)
          break LOOP //  Interrupted
        }
        command := msg[2]

        switch command {
        case "routes":
          r := show_routes()
          _ ,err = admin.SendMessage("admin", "", r); if err != nil {
            fmt.Println("err", err)
          }          
        }
        
      
      // Communications with the responders  
      case backend: 
        msg, err := s.RecvMessage(0)
        if err != nil {
          danger("Backend", "Cannot receive message", err)
          break LOOP //  Interrupted
        }

        identity, routeid, msg := unwrap_responder_message(msg)
        
        if msg != nil {
          _, err := frontend.SendMessage(routeid, "", msg[0], msg[1], msg[2])
          if err != nil {
            danger("Backend", "Cannot send message", err)
          } else {
            // find the parked responder and move it back to the queue
            i := index_of(parked[routeid], identity)
            parked[routeid] = remove_from(parked[routeid], i)
            responders[routeid] = append(responders[routeid], identity)
          }
        } else {
          if responders[routeid] == nil {
            responders[routeid] = make([]string, 0)
            parked[routeid]  = make([]string, 0)
            info("Added new route:", routeid )
          }
          responders[routeid] = append(responders[routeid], identity)
          info("Added new responder", identity, "to route:", routeid)
          fmt.Println("Added new responder", identity, "to route:", routeid)
          info("There are now", len(responders[routeid]), "in the route", routeid)
        }


      // Communciations with the acceptors
      case frontend:      
        msg, err := s.RecvMessage(0)
        if err == nil {
          
          routeid, msg := unwrap_client_message(msg)
          info(routeid, msg)
          rand.Seed(time.Now().UTC().UnixNano())

          if len(responders[routeid]) > 0 {
            // send the responder identity, followed by an empty frame and
            // the message from the client
            _, err := backend.SendMessage(responders[routeid][0], "", msg)
            info("Sending message to", responders[routeid][0])
            if err != nil {
              danger("Frontend", "Cannot send message to", routeid, "because:", err)
            }  else {
              // move the responder to parked
              parked[routeid] = append(parked[routeid], responders[routeid][0])
              responders[routeid] = responders[routeid][1:]
            }      
          } else {
            info("Route not found:", routeid)
            _, err := frontend.SendMessage(routeid, "", "404", "{\"Content-Type\" : \"text/plain; charset=UTF-8\"}", "Route not found")
            if err != nil {
              danger("Frontend", "Cannot send message for route not found", err)
            } 
          }
        }
      }
    }
  }
}


// for unwrapping the client and responder messages

func unwrap_responder_message(msg []string) (identity string, routeid string, data []string) {
  identity = msg[0]
  if msg[1] == "" {
    routeid = msg[2]
    if len(msg[3:]) > 0 {
      data = msg[3:]
    }    
  }
  return
}

func unwrap_client_message(msg []string) (routeid string, data []string) {
  routeid = msg[0]
  if len(msg) > 1 && msg[1] == "" {
    data = msg[2:]
  } else {
    data = msg[1:]
  }
  return
}

// for logging

func info(args ...interface{}) {
  logger.SetPrefix("INFO ")
  logger.Println(args...)
}

func danger(args ...interface{}) {
  logger.SetPrefix("ERROR ")
  logger.Println(args...)
}

func warning(args ...interface{}) {
  logger.SetPrefix("WARNING ")
  logger.Println(args...)
}

// for manipulating the responder and parked maps

func index_of(slice []string, id string) int {
  for index, value := range slice {
    if (value == id) {
      return index
    }
  }
  return -1  
}

func remove_from(slice []string, position int) []string {
  s1 := slice[:position]
  s2 := slice[position+1:]
  return append(s1, s2...)  
}

// 

func show_routes() []string {
  fmt.Println("getting routes", responders)
  r := make([]string, len(responders))
  for k, _ := range responders {
    r = append(r, k)
  }
  
  return r
}