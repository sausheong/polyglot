package main

import (
  "os"
  "fmt"
  "github.com/codegangsta/cli"
  zmq "github.com/pebbe/zmq4"
)

func main() {
  app := cli.NewApp()
  app.Name = "cladmin"
  app.Usage = "Command line administration tool for Polyglot broker"
	client, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		panic(err)
	}
  client.SetIdentity("admin")
  client.Connect("tcp://localhost:9999")
  defer client.Close()
  
  app.Commands = []cli.Command{
    {
      Name:      "routes",
      ShortName: "r",
      Usage:     "show routes",
      Action: func(c *cli.Context) {
        client.SendMessage("routes")        
        reply, err := client.RecvMessage(0); if err != nil {
          fmt.Println("err:", err)
        } else {
          fmt.Println(reply)
        }
      },
    },
  }
  
  app.Run(os.Args)
}