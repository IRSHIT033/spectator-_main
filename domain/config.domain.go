package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigDetails struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"` 
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string             `bson:"name" json:"name" validate:"required"`
    
}

type SiteConfig struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"` 
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	SiteUrl   string             `bson:"site_url" json:"site_url" validate:"required"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Password  string             `bson:"password" json:"password" validate:"required"`
}

