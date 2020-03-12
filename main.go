package main

import (
	"blachat-server/config"
	"blachat-server/db"
	"blachat-server/server"
	"flag"
	"fmt"
	"os"
)

func main() {

	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server - e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)

	db.Init()
	db.InitRedis()
	server.Init()
}
