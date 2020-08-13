package main

import "github.com/caeril/totpd/config"
import "github.com/caeril/totpd/server"
import "github.com/caeril/totpd/data"

func main() {

	config.InitConfig()
	data.InitData()
	server.InitTemplates()
	server.InitHandlers()
	server.InitRoutes()
	server.Run()

}
