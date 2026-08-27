package cache

import (
	"HailowSellerService/internal/domain"
	"HailowSellerService/pkg/logging"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	logger *logging.Logger
}

func New(host string, port string, logger *logging.Logger) (*RedisClient, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr: address,

		DialTimeout: 5 * time.Second,

		ReadTimeout: 3 * time.Second,

		WriteTimeout: 3 * time.Second,

		PoolTimeout: 4 * time.Second,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		Client: client,
		logger: logger,
	}, nil
}

func (client *RedisClient) CloseClient() {
	err := client.Close()
	if err != nil {
		client.logger.Errorf("%s: %v", domain.ErrRedisClose.Message, err)
	}
	client.logger.Info("Redis client closed")
}
