package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"spectator.main/domain"
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
// 		// Save the result to the database
// 		checkLog := domain.CheckLog{
// 			UrlRegionID:    1,
// 			ResponseTime:   latency,
// 			ResponseStatus: status,
// 		}
		

// 	})
// 	c.Start()
// 	select {}
// }
