package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"spectator.main/domain"
	"spectator.main/rabbitmq"

)

type CheckLogRepository struct {
	checkLogRepo   domain.CheckLogRepository
	contextTimeout time.Duration
}

// InsertOne implements domain.CheckLogUsecase.
func (c *CheckLogRepository) InsertOne(ctx context.Context, u *domain.CheckLog) (*domain.CheckLog, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()

	res, err := c.checkLogRepo.InsertOne(ctx, u)
	if err != nil {
		return res, err
	}

	return res, nil
}

func pingURL(url string) (Duration time.Duration, status any, error error) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return 0, 400, err
	}

	latency := time.Since(start)
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	message := fmt.Sprintf("Ping to %s took %v, status: %d", url, duration, status)
	routingKey := "ping/server_name"

	err = publishMessage(conn, routingKey, message)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}else {
	fmt.Println("Ping failed, no message published.")
}
	return latency, resp.StatusCode, nil
}
func NewCheckLog(u domain.CheckLogRepository, to time.Duration) domain.CheckLogUsecase {
	return &CheckLogRepository{
		checkLogRepo:   u,
		contextTimeout: to,
	}
}

// func NewCheckLogUsecase() {
// 	c := cron.New()
// 	c.AddFunc("* * * * *", func() {
// 		latency, status, err := pingURL("https://example.com/")
// 		if err != nil {
// 			return
// 		}
// 		
// 		checkLog := domain.CheckLog{
// 			UrlRegionID:    1,
// 			ResponseTime:   latency,
// 			ResponseStatus: status,
// 		}
		

// 	})
// 	c.Start()
// 	select {}
// }
