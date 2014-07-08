package main

import (
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/_/*parameters", input, process, output)
  r.POST("/_/*parameters", input, process, output)
  r.Run(":8080")
}

