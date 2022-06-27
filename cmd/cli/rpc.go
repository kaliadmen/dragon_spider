package main

import (
	"github.com/fatih/color"
	"log"
	"net/rpc"
	"os"
)

func rpcClient(inMaintenance bool) {
	port := os.Getenv("RPC_PORT")
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		gracefulExit(err)
	}
	var result string

	log.Println("Connected...")
	err = client.Call("RPCServer.MaintenanceMode", inMaintenance, &result)
	if err != nil {
		gracefulExit(err)
	}

	color.Yellow(result)
}
