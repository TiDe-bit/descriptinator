package main

import "descriptinator/pkg/server"

func main() {
	servitor := server.NewServinator("localhost:8080")
	servitor.Serve()
}
