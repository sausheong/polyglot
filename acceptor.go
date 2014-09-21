package main

import (
  "fmt"
  "net/http"
  "time"
  "github.com/julienschmidt/httprouter"  
)

func main() {
  router := httprouter.New()
  
  router.GET("/_/*p", process)
  router.POST("/_/*p", process)
  
  router.ServeFiles("/_static/*filepath", http.Dir(config.Static))
  
  server := &http.Server{
    Addr:           config.Acceptor,
    Handler:        router,
    ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
    WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
    MaxHeaderBytes: 1 << 20,
  }
  fmt.Println("Polyglot Acceptor", version(), "started at", config.Acceptor)
  server.ListenAndServe()  
}

