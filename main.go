package main

import "awesome-bluebook/ioc"

func main() {
	server := ioc.InitWebServer()
	server.Run(":8080")
}
