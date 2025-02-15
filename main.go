package main

import (
	"fmt"
	"log"
	"net/http"
	"openai/bootstrap"
	"openai/config"
	"openai/internal/handler"
	"os"
)

// GO版本如果不一样，直接改go.mod即可
func main() {
	r := bootstrap.New()

	// 微信消息处理
	r.POST("/wx", handler.ReceiveMsg)
	// 用于公众号自动验证
	r.GET("/wx", handler.WechatCheck)
	// 用于测试 curl "http://127.0.0.1:$PORT/test"
	r.GET("/test", handler.Test)

	// 设置日志
	SetLog()

	fmt.Printf("启动服务，使用 curl 'http://127.0.0.1:%s/test?msg=中国在哪个洲' 测试一下吧\n", config.ServerPort)
	if err := http.ListenAndServe(":"+config.ServerPort, r); err != nil {
		panic(err)
	}
}

func SetLog() {
	dir := "./log"
	file := dir + "/chatgpt.log"
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	fmt.Println("查看日志请使用 tail -f " + file)
}
