package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"spectator.main/internals/bootstrap"
	_userRepo "spectator.main/user/repository/mongo_repository"
	_userUsecase "spectator.main/user/usecase"
	_userHandler "spectator.main/user/transport/http"
)

func main() {
	router := gin.Default()
	
	gin.SetMode(gin.DebugMode)

    timeoutContext := time.Duration(bootstrap.App().Config.ContextTimeout) * time.Second
    
	database:= bootstrap.App().Mongo.Database(bootstrap.App().Config.DBName)

    app:=router.Group("api/v1")
    
	userRepo:= _userRepo.NewMongoRepository(database)
	userUseCase:= _userUsecase.NewUserUsecase(userRepo,timeoutContext)
	_userHandler.NewUserHandler(app,userUseCase)
	
	router.Run(":8080")
}
