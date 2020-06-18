package main

import (
	"fmt"
	"go_nat_git/Service"
)

func main() {
	var clientConfig = Service.GetClientConfig()
	fmt.Println("client is", clientConfig)
	var serverConfig = Service.GetServerConfig()
	fmt.Println("server is", serverConfig)
}
