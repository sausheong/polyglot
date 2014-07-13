package main

import (
  "github.com/nu7hatch/gouuid"
  "github.com/streadway/amqp"

  "log"
  "fmt"
  "strings"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}


func main() {
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()

  routeId = "GET/_/go/hello"
  
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()
    
  // declare the response queue used to receive responses from the responders
  queue, err := ch.QueueDeclare(
    routeId,            // name
    true,              // durable
    true,               // delete when unused
    false,              // exclusive
    false,              // noWait
    nil,                // arguments
  )    
  
  println("[Responder ready.]")
  
  // wait to receive 
  msgs, err := ch.Consume(
    queue.Name,     // queue
    "",             // consumer
    true,            // auto acknowledge
    false,           // exclusive
    false,           // no local
    false,           // no wait
    nil,             // table
  )
  failOnError(err, "Failed to consume message")
  
  ret := make(chan []byte)
  go func() {
    for d := range msgs {
      ret <- d.Body
    }
  }()  
  response := string(<-ret)  
 
 
  
}

