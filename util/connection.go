package util

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

// Connection struct to pass into the service
type Connection struct {
	Backend  string
	Postgres *gorm.DB
	Redis    *redis.Client
}
