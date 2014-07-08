package main

import "github.com/gin-gonic/gin"

func main() {
  r := gin.Default()
  
  r.GET("/_/*p", process)
  r.POST("/_/*p", process)
  
  r.Run(":8080")
}

