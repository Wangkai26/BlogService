package global

import (
	"github.com/go-programming-tour-book/blog-service/pkg/logger"
	"github.com/go-programming-tour-book/blog-service/pkg/setting"
)

var(
	ServerSetting	*setting.ServerSettingS
	AppSetting		*setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS

	// 在全局变量中新增 Logger对象，用于日志组件的初始化
	Logger 			*logger.Logger

	// 8section,新增JWT配置的全局对象
	JWTSetting      *setting.JWTSettingS
)
