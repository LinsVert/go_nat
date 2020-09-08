package main

import (
	"fmt"
	"go_nat_git/Service"
	"go_nat_git/client"
	"os"
)

func main() {
	daemonFlagName := "--daemon"
	isDaemon := true
	for i := 1; i < len(os.Args); i++ {
		fmt.Println(os.Args[i])
		if os.Args[i] == daemonFlagName {
			isDaemon = false
		}
	}
	if isDaemon {
		pid, err := Service.Daemon()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Println(pid, "run daemon on pid")
		os.Exit(0)
	}
	client.Run()
}
