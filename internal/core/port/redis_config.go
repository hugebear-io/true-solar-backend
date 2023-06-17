package port

import "time"

type RedisConfigRepoPort interface {
	GetRedisConfig() (RedisConfig, error)
	UpdateRedisConfig(redisConfig RedisConfig) error
}

type RedisConfig struct {
	ID        int       `json:"id"`
	Hostname  string    `json:"hostname" binding:"required"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
