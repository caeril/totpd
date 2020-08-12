package main

import "totpd/server"
import "totpd/data"

func main() {

	data.InitData()
	server.InitTemplates()
	server.InitHandlers()
	server.InitRoutes()
	server.Run()

}
