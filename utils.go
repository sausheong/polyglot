package main

import (
  "fmt"
  "time"
  "log"
  "crypto/rand"
  "encoding/json"
  "os"  
)

type Configuration struct {
  Acceptor        string            
  ReadTimeout     int64
  WriteTimeout    int64  
  RequestTimeout  int64
  RequestRetries  int8
  Broker          string
}

func (config *Configuration) Timeout() time.Duration {
  return time.Duration(config.RequestTimeout * int64(time.Millisecond))
}

var config Configuration
var logger *log.Logger

func init() {
  loadConfig()
  file, err := os.OpenFile("acceptor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
      log.Fatalln("Failed to open log file", err)
  }  
  logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)  
}

func loadConfig() {
  file, err := os.Open("config.json")
  failOnError(err, "Cannot open config file")
  decoder := json.NewDecoder(file)
  config = Configuration{}
  err = decoder.Decode(&config)
  failOnError(err, "Cannot get configuration from file")  
}


// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
  u := new([16]byte)
  _, err := rand.Read(u[:])
  failOnError(err, "Cannot generate UUID")
  // 0x40 is reserved variant from RFC 4122  
  u[8] = (u[8] | 0x40) & 0x7F
  // Set the four most significant bits (bits 12 through 15) of the
  // time_hi_and_version field to the 4-bit version number.  
  u[6] = (u[6] & 0xF) | (0x4 << 4)
  uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
  return
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
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