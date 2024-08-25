package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"spectator.main/internals/bootstrap"
	rabbitm "spectator.main/ping/rabbitmq"
	_checkLogRepo "spectator.main/ping/repository/mongo_repository"
	_checkLogUsecase "spectator.main/ping/usecase"
	_userRepo "spectator.main/user/repository/mongo_repository"
	_userHandler "spectator.main/user/transport/http"
	_userUsecase "spectator.main/user/usecase"
)
func schedulePingJob() {
	
	c := cron.New()

	   _, err := c.AddFunc("@every 4m", func() {
		url := "http://example.com" 
		duration, status, err := _checkLogUsecase.PingURL(url)
		if err != nil {
			
			fmt.Printf("Error pinging URL: %v\n", err)
			return
		}

		fmt.Printf("Pinged %s in %v, status: %d\n", url, duration, status)

		conn, err := rabbitmq.ConnectToRabbitMQ()
		if err != nil {
			fmt.Printf("Failed to connect to RabbitMQ: %v\n", err)
			return
		}
		defer conn.Close()

		if status == 200 {
			message := fmt.Sprintf("Ping to %s took %v, status: %d", url, duration, status)
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

	timeoutContext := time.Duration(config.ContextTimeout) * time.Second

	database := app.Mongo.Database(config.DBname)

	ginRouter := router.Group("api/v1")

	userRepo := _userRepo.NewMongoRepository(database)
	userUseCase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(ginRouter, userUseCase)
	 _checkLogRepo.NewMongoRepository(database)
	

	router.Run(":8080")
}
