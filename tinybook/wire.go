//go:build wireinject

package main

import (
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/service/sms/localsms"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {

	wire.Build(
		// 初始化redis 和 db
		ioc.InitRedis, ioc.InitDB,
		// 初始化user模块
		cache.NewUserCache, dao.NewUserDAO, repository.NewUserRepository, service.NewUserService,
		// 初始化code模块
		cache.NewCodeCache, repository.NewCodeRepository, service.NewCodeService, localsms.NewService,
		// 初始化handler
		web.NewUserHandler,
		// 初始化web
		ioc.InitWebServer, ioc.InitHandlerFunc,
	)

	return gin.Default()
}
