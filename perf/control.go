package main

import (
  "fmt"
  "time"
  "net/http"
)

func work(w http.ResponseWriter, r *http.Request) {
  time.Sleep(500 * time.Millisecond)
  fmt.Fprintf(w, "<h1>Hello Perf</h1>")
}

func main() {
  http.HandleFunc("/perf", work)
  http.ListenAndServe("0.0.0.0:8080", nil)
}

