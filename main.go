package main

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
import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"takeout/global"
	"takeout/initialize"
	"time"
)

func main() {

	router := initialize.GlobalInit()
	gin.SetMode(global.Config.Server.Level)
	fmt.Println("[SERVER] runs on http://localhost:8080")

	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + global.Config.Server.Port,
		Handler: router,
	}

	// 协程运行并监听服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("an error occurred while running the server: %s\n", err)
		}
	}()

	// 启动 pprof 监听
	go func() {
		log.Println("[PPROF] Running on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("failed to start pprof server: %s\n", err)
		}
	}()

	// 退出系统
	quit := make(chan os.Signal, 1)
	// 创建缓存为1的信号量channel
	signal.Notify(quit, os.Interrupt)
	// 阻塞等待取出信号量
	<-quit
	// 接收到信号量后打印提示
	log.Println("server is being shutdown...")
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 延迟调用cancel，确保上下文能取消
	defer cancel()
	// 调用Shutdown关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error:", err)
	}
	log.Println("server has been safely shutdown")

}
