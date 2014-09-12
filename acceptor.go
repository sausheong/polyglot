package main

import (
  "fmt"
  "net/http"
  "time"
  "github.com/julienschmidt/httprouter"  
)

func main() {
  fmt.Println("Polyglot Acceptor v0.1")
  fmt.Println("")
  router := httprouter.New()
  
  router.GET("/_/*p", process)
  router.POST("/_/*p", process)
  
  
  server := &http.Server{
    Addr:           config.Acceptor,
    Handler:        router,
    ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
    WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
    MaxHeaderBytes: 1 << 20,
  }
  server.ListenAndServe()  

}

