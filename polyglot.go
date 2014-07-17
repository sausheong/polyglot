package main

import (
  "github.com/julienschmidt/httprouter"
  "github.com/nu7hatch/gouuid"
  "github.com/streadway/amqp"
  "encoding/json"
  "encoding/base64"
  "log"
  "fmt"
  "strings"
  "net/http"
  // "reflect"
)


func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}


// default handler
func process(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {  
  r.ParseForm()
  r.ParseMultipartForm(1024)

  // marshal the HTTP request struct into JSON
  req_json, err := json.Marshal(r)
  
  routeId := r.Method + r.URL.Path
  fmt.Println(routeId)
  failOnError(err, "Failed to marshal the request")  
  
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()

  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()

  _, err = ch.QueueInspect(routeId); if err != nil {
    w.WriteHeader(404)
    w.Write([]byte("Not Found"))
    fmt.Println(" - Not Found")
    return
  }

  // declare the response queue used to receive responses from the responders
  replyq, err := ch.QueueDeclare(
    routeId + ":[r]",    // name
    false,              // durable
    false,               // delete when unused
    false,              // exclusive
    false,              // noWait
    nil,                // arguments
  )
  
  // assert type of the body
  body := req_json
  
  // publish the request into the polyglot queue
  corrId, _ := uuid.NewV4()

  
  err = ch.Publish(
    "",         // default exchange
    routeId,    // routing key
    false,      // mandatory
    false,
    amqp.Publishing {
      DeliveryMode:  amqp.Persistent,
      ContentType:   "application/json",
      CorrelationId: corrId.String(),
      ReplyTo:       replyq.Name,
      Body:          []byte(body),
      AppId:         routeId,
    })
  failOnError(err, "Failed to publish a message")  
  fmt.Println(" - published to queue")
  // wait to receive 
  msgs, err := ch.Consume(
    replyq.Name,     // queue
    "process",       // consumer
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
  err = ch.Cancel("process", false)
  failOnError(err, "Failed to cancel channel")   
  
  // get response JSON array 
  res := string(response)

  // unmarshal JSON into status, headers and body
  var resp interface{}
  err = json.Unmarshal([]byte(res), &resp); if err == nil {
    response := resp.([]interface{})
    status := response[0]
    headers := response[1].(map[string]interface{})
    body := response[2]

    // write headers
    for k, v := range headers {
      w.Header().Set(k, v.(string))
    }
    s, _ := status.(float64)
    b, _ := body.(string)
    var data []byte
    
    // get content type    
    ctype, hasCType := headers["Content-Type"].(string); if hasCType == true {
      if strings.HasPrefix(ctype, "text") {
        data = []byte(b)
      } else {
        data, _ = base64.StdEncoding.DecodeString(b)
      }
    } else {
      data, _ = base64.StdEncoding.DecodeString(b)
    }
    fmt.Println(" -", status)
    // write status and body to response
    w.WriteHeader(int(s))
    w.Write(data)
  }
}