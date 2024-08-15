package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"spectator.main/internals/bootstrap"
   _checkLogRepo "spectator.main/ping/repository/mongo_repository"
   _checkLogUsecase "spectator.main/ping/usecase"
	_userRepo "spectator.main/user/repository/mongo_repository"
	_userHandler "spectator.main/user/transport/http"
	_userUsecase "spectator.main/user/usecase"
)

func main() {

	app := bootstrap.App()

	config := app.Config

	router := gin.Default()

	gin.SetMode(gin.DebugMode)

	timeoutContext := time.Duration(config.ContextTimeout) * time.Second

	database := app.Mongo.Database(config.DBname)

	ginRouter := router.Group("api/v1")

	userRepo := _userRepo.NewMongoRepository(database)
	userUseCase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(ginRouter, userUseCase)
	 _checkLogRepo.NewMongoRepository(database)
	_checkLogUsecase.NewCheckLogUsecase()

	router.Run(":8080")
}
