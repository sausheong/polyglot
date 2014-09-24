package main

import (  
  "encoding/json"
  "encoding/base64"
  "fmt"
  "strings"
  "strconv"
  "net/http"
  
  "github.com/julienschmidt/httprouter"
  zmq "github.com/pebbe/zmq4"
)


// default handler
func process(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  if request.Method == "POST" {
    err := request.ParseForm()
    failOnError(err, "Failed to parse form")
    request.ParseMultipartForm(10485760)
  }
  
  // marshal the HTTP request struct into JSON
  reqJson, err := json.Marshal(request)
  failOnError(err, "Failed to marshal the request into JSON")
  
  routeId := request.Method + request.URL.Path

  // send request
  // create ZMQ REQ socket
	client, err := zmq.NewSocket(zmq.REQ)
  failOnError(err, "Failed to create socket")
  defer client.Close()
  
  // set the identity to the route ID eg GET/_/path
  client.SetIdentity(routeId)
	client.Connect(config.Broker)

	poller := zmq.NewPoller()
	poller.Add(client, zmq.POLLIN)

	retries_left := config.RequestRetries
  RETRIES_LOOP:
	for retries_left > 0 {
		//  We send a request, then we work to get a reply
		client.SendMessage(reqJson)

    
		for expect_reply := true; expect_reply; {
			//  Poll socket for a reply, with timeout
			sockets, err := poller.Poll(config.Timeout()); if err != nil {
          reply(writer, 500, []byte(err.Error()))
          retries_left = 0
          break RETRIES_LOOP
        }
      
      // if there is a reply
			if len(sockets) > 0 {
				response, err := client.RecvMessage(0); if err != nil {
          reply(writer, 500, []byte(err.Error()))
          retries_left = 0
          break RETRIES_LOOP
				}
				info(response[0])

        status := response[0] // HTTP response code eg 200, 404
        headers_json := response[1] // JSON encoded HTTP response headers
        body := response[2] // HTTP response body as a string

        // // unmarshal header JSON
        var headers map[string]string
        err = json.Unmarshal([]byte(headers_json), &headers); if err != nil {
          reply(writer, 500, []byte(err.Error()))
          retries_left = 0
          break RETRIES_LOOP
        }
 
        // write headers
        for k, v := range headers {
          writer.Header().Set(k, v)
        }
        s, err := strconv.Atoi(status); if err != nil {
          reply(writer, 500, []byte(err.Error()))
          retries_left = 0
          break RETRIES_LOOP
        }
        

        var data []byte
        // get content type
        ctype, hasCType := headers["Content-Type"]; if hasCType == true {
          if is_text_mime_type(ctype) {
            data = []byte(body)
          } else {
            data, _ = base64.StdEncoding.DecodeString(body)
          }
        } else {
          data = []byte(body) // if not given the content type, assume it's text
        }

        // write status and body to response
        reply(writer, s, data)

				retries_left = 0
				expect_reply = false

      // if there are no replies, try again  
			} else {
				retries_left--
				if retries_left == 0 {
          reply(writer, 500, []byte("Cannot connect with broker, giving up."))
					break
				} else {
					fmt.Println("Cannot reach broker, retrying...")
					//  Old socket is confused; close it and open a new one
					client.Close()
					client, err = zmq.NewSocket(zmq.REQ)
          failOnError(err, "Failed to create socket")
          defer client.Close()
          client.SetIdentity(routeId)
					client.Connect(config.Broker)
					// Recreate poller for new client
					poller = zmq.NewPoller()
					poller.Add(client, zmq.POLLIN)
					//  Send request again, on new socket
					client.SendMessage(reqJson)
				}
			}
		}
	}
}

func reply(writer http.ResponseWriter, status int, body []byte) {
  writer.WriteHeader(status)
  writer.Write(body)
}

func is_text_mime_type(ctype string) bool {
  if strings.HasPrefix(ctype, "text") ||
  strings.HasPrefix(ctype, "application/json") {
    return true
  } else {
    return false
  }
  
}
