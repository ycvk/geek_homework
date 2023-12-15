// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"geek_homework/tinybook/internal/events/article"
	"geek_homework/tinybook/internal/events/interactive"
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/repository/cache"
	"geek_homework/tinybook/internal/repository/dao"
	"geek_homework/tinybook/internal/service"
	"geek_homework/tinybook/internal/web"
	"geek_homework/tinybook/internal/web/jwt"
	"geek_homework/tinybook/ioc"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := ioc.InitLogger()
	v := ioc.InitHandlerFunc(cmdable, handler, logger)
	db := ioc.InitDB(logger)
	userDAO := dao.NewGormUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository, logger)
	theineCache := ioc.InitLocalCache()
	codeCache := cache.NewLocalCodeCache(theineCache)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	smsdao := dao.NewGormSMSDAO(db)
	smsRepository := repository.NewGormSMSRepository(smsdao)
	smsService := ioc.InitSMSService(cmdable, smsRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	database := ioc.InitMongoDB(logger)
	conn := ioc.InitMongoDBV2()
	articleDAO := dao.NewMongoDBArticleDAO(database, conn)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, articleCache, userRepository, logger)
	writer := ioc.InitWriter()
	readEventProducer := article.NewKafkaArticleProducer(writer)
	articleService := service.NewArticleService(articleRepository, readEventProducer, logger)
	interactiveDAO := dao.NewGormInteractiveDAO(db)
	likeRankEventProducer := interactive.NewKafkaLikeRankProducer(writer)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable, logger, theineCache, likeRankEventProducer)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository, articleRepository, likeRankEventProducer, logger)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, logger)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	kafkaConsumer := article.NewKafkaConsumer(interactiveRepository, logger)
	interactiveKafkaConsumer := interactive.NewKafkaLikeRankConsumer(logger, theineCache, cmdable)
	v2 := article.CollectConsumer(kafkaConsumer, interactiveKafkaConsumer)
	app := &App{
		server:    engine,
		consumers: v2,
	}
	return app
}
