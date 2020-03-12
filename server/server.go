package server

import "blachat-server/config"

func Init(){
	config := config.GetConfig()
	r := NewRouter()
	_ = r.Run(config.GetString("server_port"))
}