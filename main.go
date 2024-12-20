package main

import (
	"fmt"
	"takeout/initialize"
)

func main() {
	router := initialize.GlobalInit()
	//gin.SetMode(global.Config.Server.Level)
	fmt.Println("[SERVER] runs on http://localhost:8080")
	router.Run(":8080")
}
