package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckLog struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UrlRegionID    int                `bson:"url_region_id" json:"url_region_id"`
	ResponseTime   time.Duration      `bson:"response_time" json:"response_time"`
	ResponseStatus int                `bson:"response_status" json:"response_status"`
}

type Response struct {
	Latency time.Duration `bson:"latency" json:"latency"`
	Status  int           `bson:"status" json:"status"`
}

type CheckLogRepository interface {
	InsertOne(ctx context.Context, u *CheckLog) (*CheckLog, error)
}

type CheckLogUsecase interface {
	InsertOne(ctx context.Context, u *CheckLog) (*CheckLog, error)
	PingURL(ctx context.Context) (*Response, error)
}
