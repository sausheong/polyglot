package main

import (
  "github.com/streadway/amqp"
  "log"
  "fmt"
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

  routeId := "GET/_/go/hello"
  
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()
  
  
  err = ch.Qos(1, 0, false)  
  failOnError(err, "Failed to set QoS")
  // declare the response queue used to receive responses from the responders
  queue, err := ch.QueueDeclare(
    routeId,            // name
    false,              // durable
    true,               // delete when unused
    false,              // exclusive
    false,              // noWait
    nil,                // arguments
  )    
  
  println("[Responder ready.]")
  
  // ret := make(chan []byte)
  cor := make(chan string)
  rep := make(chan string)
  aid := make(chan string)
    
  // wait to receive 
  for {
    deliveries, err := ch.Consume(
      queue.Name,     // queue
      "",             // consumer
      false,            // auto acknowledge
      false,           // exclusive
      false,           // no local
      false,           // no wait
      nil,             // table
    )
    failOnError(err, "Failed to consume message")

    go func() {
      for d := range deliveries {
        d.Ack(true)
        // ret <- d.Body
        cor <- d.CorrelationId
        rep <- d.ReplyTo
        aid <- d.AppId
      }
    }()  
    // response := string(<-ret) 
    corrId := string(<-cor)
    replyTo := string(<-rep)  
    appId := string(<-aid)
    
    if routeId == appId {
      err = ch.Publish(
        "",         // default exchange
        replyTo,    // routing key
        false,      // mandatory
        false,
        amqp.Publishing {
          DeliveryMode:  amqp.Persistent,
          ContentType:   "text/html",
          CorrelationId: corrId,
          Body:          []byte("[200, {\"Content-Type\": \"text/html\"}, \"<h1>Hello Go!</h1>\"]"),
          AppId:         routeId,
        })
      failOnError(err, "Failed to publish a message")    
    
      // fmt.Println("Response:", response)
     
      err = ch.Cancel("process", true)
      failOnError(err, "Failed to cancel channel")          
    }
  }
}

