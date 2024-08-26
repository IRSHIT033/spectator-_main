package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"spectator.main/domain"
	"spectator.main/internals/bootstrap"
	_checkLogRepo "spectator.main/ping/repository/mongo_repository"
	rabbitmq "spectator.main/ping/transport/rabbitmq"
	"spectator.main/ping/usecase"
	_userRepo "spectator.main/user/repository/mongo_repository"
	_userHandler "spectator.main/user/transport/http"
	_userUsecase "spectator.main/user/usecase"
)
func schedulePingJob(_checkLogUsecase domain.CheckLogUsecase) {
	c := cron.New()

	err := c.AddFunc("1 * * * * *", func() {
		ctx := context.Background()


		response, err :=  _checkLogUsecase.PingURL(ctx)
		if err != nil {
			fmt.Printf("Error pinging URL: %v\n", err)
			return
		}
		duration := response.Latency
		status := response.Status

		fmt.Printf("Pinged URL in %v, status: %d\n", duration, status)

		conn, err := rabbitmq.ConnectToRabbitMQ()
		if err != nil {
			fmt.Printf("Failed to connect to RabbitMQ: %v\n", err)
			return
		}
		defer conn.Close()

		if status == 200 {
			message := fmt.Sprintf("Ping took %v, status: %d", duration, status)
			routingKey := "ping/server_name"

			err = rabbitmq.PublishMessage(conn, routingKey, message)
			if err != nil {
				fmt.Printf("Failed to publish message: %v\n", err)
			} else {
				fmt.Println("Message published successfully.")
			}
		} else {
			fmt.Println("Ping failed, no message published.")
		}
	})

	if err != nil {
		fmt.Printf("Failed to schedule ping job: %v\n", err)
		return
	}

	c.Start()
}





func main() {

	app := bootstrap.App()

	config := app.Config

	router := gin.Default()

	gin.SetMode(gin.DebugMode)
   
	timeoutContext := time.Duration(12) * time.Second
     fmt.Println("TIMEOUT",timeoutContext)
	database := app.Mongo.Database(config.DBname)

	ginRouter := router.Group("api/v1")

	userRepo := _userRepo.NewMongoRepository(database)
	userUseCase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(ginRouter, userUseCase)
	mongoRepo:=_checkLogRepo.NewMongoRepository(database)
	_newLogUsecase := usecase.NewCheckLog(mongoRepo, 10*time.Second)

	schedulePingJob(_newLogUsecase)

	router.Run(":8080")
}
