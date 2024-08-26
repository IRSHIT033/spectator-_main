package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	fmt.Print("RESPONSEEEEEEEEEEEEEEEEEEEEEEEEEE",res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// PingURL pings the given URL and inserts the latency and status into the database.
func (c *CheckLogRepository) PingURL(ctx context.Context) (*domain.Response, error) {
	
	start := time.Now()
	url := "http://example.com"
	urlRegionID := 1

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	latency := time.Since(start)
 
	checkLog := &domain.CheckLog{
		UrlRegionID:    urlRegionID,
		ResponseTime:   latency,
		ResponseStatus: resp.StatusCode,
	}
	fmt.Printf("STARTEDDDDDDDDDDDDDDDDDDDDDDDDD");
	_, err = c.InsertOne(ctx, checkLog)

	fmt.Printf("INSERTEDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD");
	if err != nil {
		fmt.Printf("ERRRORRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR");
		return nil, err
	}

	return &domain.Response{
		Latency: latency,
		Status:  resp.StatusCode,
	}, nil
}

func NewCheckLog(u domain.CheckLogRepository, to time.Duration) domain.CheckLogUsecase {
	return &CheckLogRepository{
		checkLogRepo:   u,
		contextTimeout: to,
	}
}
