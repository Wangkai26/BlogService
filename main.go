package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/internal/routers"
	"github.com/go-programming-tour-book/blog-service/pkg/logger"
	"github.com/go-programming-tour-book/blog-service/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)

func init()  {
	err := setupSetting()
	if err != nil{
		log.Fatalf("init.setupSetting err: %v",err)
	}
	err = setupDBEngine()
	if err != nil{
		log.Fatalf("init.setupDBEngine err: %v",err)
	}
	err = setupLogger()
	if err != nil{
		log.Fatalf("init.setupLogger err: %v",err)
	}
}

// @title 博客系统
// @version 1.0
// @description Go 语言编程之旅：一起用Go做项目
// @termsOfService https://github.com/go-programming-tour-book
func main()  {
	gin.SetMode(global.ServerSetting.RunMode)
	// 下面三行打印用来校验配置是否真正地映射到了配置结构体上
	//log.Println(global.AppSetting)
	//log.Println(global.ServerSetting)
	//log.Println(global.DatabaseSetting)
	router := routers.NewRouter()
	s := &http.Server{
		Addr: 			":"+global.ServerSetting.HttpPort,
		Handler: 		router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1<<20,
	}
	global.Logger.Infof("%s:go-programming-tour-book/%s","eddycjy","blog-service")
	s.ListenAndServe()
}

func setupSetting() error {
	setting,err := setting.NewSetting()
	if err != nil{
		return err
	}
	err = setting.ReadSection("Server",&global.ServerSetting)
	if err != nil{
		return err
	}
	err = setting.ReadSection("App",&global.AppSetting)
	if err != nil{
		return err
	}
	err = setting.ReadSection("Database",&global.DatabaseSetting)
	if err != nil{
		return err
	}
	// 8section,设置JWT的一些相关配置
	err = setting.ReadSection("JWT",&global.JWTSetting)
	if err != nil{
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
	return nil
}

func setupDBEngine() error {
	var err error
	// 这里一定不能用 :=
	// 由于 := 会重新声明并创建左侧的新局部变量，会导致等式右方变量没有赋值给全局变量 global.DBEngine
	// 其他包调用 global.DBEngine 时，它仍然是 nil
	global.DBEngine,err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil{
		return err
	}

	return nil
}

func setupLogger() error {
	filename := global.AppSetting.LogSavePath+"/"+global.AppSetting.LogFileName+global.AppSetting.LogFileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: filename,
		MaxSize:  600,
		MaxAge:   10,
		LocalTime: true,
	},"",log.LstdFlags).WithCaller(2)

	return nil
}



