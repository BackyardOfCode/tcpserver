package main
import (
	"log"
	"os"
	"tcpserver/server"
)

func main() {

	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	log.SetOutput(file) // 设置输出目标

	// 自定义日志前缀（日期时间 + 自定义标记）
	log.SetPrefix("[MyApp] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 控制日志格式
	log.Println("bbb,,,,nnnnn")
	return
	s := server.NewServer("0.0.0.0", 8099)
	s.Start()
}
