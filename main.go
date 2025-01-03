package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"takeout/global"
	"takeout/initialize"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /
//
// @securityDefinitions.apikey JWTAuth
// @in header
// @name token

func main() {

	router := initialize.GlobalInit()
	gin.SetMode(global.Config.Server.Level)
	fmt.Println("[SERVER] runs on http://localhost:8080")
	router.Run(":8080")

}
